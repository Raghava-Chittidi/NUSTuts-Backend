package dataaccess

import (
	"NUSTuts-Backend/internal/database"
	"NUSTuts-Backend/internal/models"
	"NUSTuts-Backend/internal/util"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/samber/lo"
)

// Created a tutorial for every single module available in NUS
func CreateTutorialsForEveryModule() error {
	// Send a request to get the list of modules for this current Academic Year
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

	/* 
		Filter such that each module's list of semesters only contains the current semester or 
		nothing if that module is not taught during the current semester 
	*/
	modulesList = lo.Map(modulesList, func(item Module, index int) Module {
		item.Semesters = lo.Filter(item.Semesters, func(sem int, index int) bool {
			return sem == util.GetCurrentSem()
		})

		return item
	})

	// Filter out the modules that will not be taught during the current semester
	modulesList = lo.Filter(modulesList, func(item Module, index int) bool {
		return len(modulesList[index].Semesters) > 0
	})


	for i, module := range modulesList {
		// Get the module information for each module in the list
		res1, err := http.Get(fmt.Sprintf("https://api.nusmods.com/v2/%s/modules/%s.json", util.GetCurrentAY(), module.ModuleCode))
		if err != nil {
			fmt.Printf("Failed to fetch module - %s info", module.ModuleCode)
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

		// Double check to see if that module is taught in the current semester
		curSem, found := lo.Find(moduleInfo.SemesterData, func(item SemesterData) bool {
			return item.Semester == util.GetCurrentSem()
		})

		if !found {
			continue
		}
		
		// Get the list of all the tutorials for the module
		tutorials := lo.Filter(curSem.Timetable, func(item Tutorial, index int) bool {
			return item.LessonType == "Tutorial"
		})

		if len(tutorials) == 0 {
			continue
		}

		// Create a Tutorial and a Teaching Assistant for each tutorial in the module's tutorial list 
		for j, item := range tutorials {
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

			// Assign the created Teaching Assistant to this tutorial created
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

func GetAllTutorialIds() (*[]int, error) {
	var tutorialIds []int
	result := database.DB.Table("tutorials").Select("id").Find(&tutorialIds)
	if result.Error != nil {
		return nil, result.Error
	}

	return &tutorialIds, nil
}

// Checks if the teaching assistant teaches that tutorial
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
