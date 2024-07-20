package dataaccess

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"NUSTuts-Backend/internal/util"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/samber/lo"
)

func CreateTutorialsForEveryModule() error {
	res, err := http.Get(fmt.Sprintf("https://api.nusmods.com/v2/%s/moduleList.json", util.GetCurrentAY()))
	if err != nil {
		return err
	}

	type Module struct {
		ModuleCode string `json:"moduleCode"`
		Title string `json:"title"`
		Semesters []int `json:"semesters"`
	}

	var modulesList []Module
	dec := json.NewDecoder(res.Body)
	err = dec.Decode(&modulesList)
	if err != nil {
		return err
	}

	modulesList = lo.Map(modulesList, func(item Module, index int) Module {
		item.Semesters = lo.Filter(item.Semesters, func(sem int, index int) bool {
			return sem == util.GetCurrentSem()
		})

		return item
	})

	modulesList = lo.Filter(modulesList, func(item Module, index int) bool {
		return len(modulesList[index].Semesters) > 0
	})

	// fmt.Println(len(modulesList))

	for i, module := range modulesList {
		fmt.Println(i)
		res1, err := http.Get(fmt.Sprintf("https://api.nusmods.com/v2/%s/modules/%s.json", util.GetCurrentAY(), module.ModuleCode))
		if err != nil {
			fmt.Println("err0")
			return err
		}

		type Tutorial struct {
			ClassNo    string `json:"classNo"`
			StartTime  string `json:"startTime"`
			EndTime    string `json:"endTime"`
			Weeks      []int  `json:"weeks"`
			Venue      string `json:"venue"`
			Day        string `json:"day"`
			LessonType string `json:"lessonType"`
			Size       int    `json:"size"`
			CovidZone  string `json:"covidZone"`
		}

		type SemesterData struct {
			Semester  int `json:"semester"`
			Timetable []Tutorial `json:"timetable"`
			CovidZones []string `json:"covidZones"`
		}

		type ModuleInfo struct {
			AcadYear                string    `json:"acadYear"`
			Description             string    `json:"description"`
			Title                   string    `json:"title"`
			AdditionalInformation   string    `json:"additionalInformation"`
			Department              string    `json:"department"`
			Faculty                 string    `json:"faculty"`
			Workload                []float64 `json:"workload"`
			GradingBasisDescription string    `json:"gradingBasisDescription"`
			ModuleCredit            string    `json:"moduleCredit"`
			ModuleCode              string    `json:"moduleCode"`
			SemesterData            []SemesterData `json:"semesterData"`
		}

		var moduleInfo ModuleInfo
		dec := json.NewDecoder(res1.Body)
		err = dec.Decode(&moduleInfo)
		if err != nil {
			continue
		}

		curSem, found := lo.Find(moduleInfo.SemesterData, func(item SemesterData) bool {
			return item.Semester == util.GetCurrentSem()
		})

		if !found {
			continue
		}
		
		tutorials := lo.Filter(curSem.Timetable, func(item Tutorial, index int) bool {
			return item.LessonType == "Tutorial"
		})

		if len(tutorials) == 0 {
			continue
		}

		for j, item := range tutorials {
			log.Println("inserting")

			name := fmt.Sprintf("ta%d", 200 * i + j)
			email := fmt.Sprintf("%s@gmail.com", name)
			passwordHash, err := util.GetPasswordHash(name)
			if err != nil {
				fmt.Println("err2")
				return err
			}

			teachingAssistant, err := CreateTeachingAssistant(name, email, passwordHash)
			if err != nil {
				fmt.Println("err3")
				return err
			}

			tutorial, err := CreateTutorial(item.ClassNo, module.ModuleCode, int(teachingAssistant.ID))
			if err != nil {
				fmt.Println("err4")
				return err
			}

			teachingAssistant.TutorialID = int(tutorial.ID)
			database.DB.Save(&teachingAssistant)
		}
	}

	return nil
}

func CreateTutorial(tutorialCode string, module string, teachingAssistantId int) (*models.Tutorial, error) {
	tutorial := &models.Tutorial{TutorialCode: tutorialCode, Module: module, TeachingAssistantID: teachingAssistantId}
	result := database.DB.Table("tutorials").Create(tutorial)
	return tutorial, result.Error
}

func GetTutorialById(id int) (*models.Tutorial, error) {
	var tutorial models.Tutorial
	result := database.DB.Table("tutorials").Where("id = ?", id).First(&tutorial)
	if result.Error != nil {
		return nil, result.Error
	}

	return &tutorial, nil
}

func GetTutorialByClassAndModuleCode(classNo string, moduleCode string) (*models.Tutorial, error) {
	var tutorial models.Tutorial
	result := database.DB.Table("tutorials").Where("module = ?", moduleCode).Where("tutorial_code = ?", classNo).First(&tutorial)
	if result.Error != nil {
		return nil, result.Error
	}

	return &tutorial, nil
}

func GetAllTutorialIDs() (*[]int, error) {
	var tutorialIds []int
	result := database.DB.Table("tutorials").Select("id").Find(&tutorialIds)
	if result.Error != nil {
		return nil, result.Error
	}

	return &tutorialIds, nil
}

func CheckIfStudentInTutorialById(studentId int, tutorialId int) (bool, error) {
	var registry models.Registry
	result := database.DB.Table("registries").Where("tutorial_id = ?", tutorialId).Where("student_id = ?", studentId).First(&registry)
	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

func CheckIfTeachingAssistantInTutorialById(teachingAssistantId int, tutorialId int) (bool, error) {
	var tutorial models.Tutorial
	result := database.DB.Table("tutorials").Where("id = ?", tutorialId).Where("teaching_assistant_id = ?", teachingAssistantId).First(&tutorial)
	if result.Error != nil {
		return false, result.Error
	}

	return true, nil
}

func DeleteTutorialById(id int) error {
	tutorial, err := GetTutorialById(id)
	if err != nil {
		return err
	}

	result := database.DB.Unscoped().Table("tutorials").Delete(&tutorial)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func DeleteTutorialByClassAndModuleCode(classNo string, moduleCode string) error {
	tutorial, err := GetTutorialByClassAndModuleCode(classNo, moduleCode)
	if err != nil {
		return err
	}

	result := database.DB.Unscoped().Table("tutorials").Delete(&tutorial)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
