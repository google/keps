package sections

import (
	_ "github.com/calebamiles/keps/pkg/keps/sections/internal/rendering"
)

func ForProvisionalState(renderingInfo renderingInfoProvider) ([]Entry, error) {
	/*
		summary, err := New(rendering.SummaryName, renderingInfo)
		if err != nil {
			return nil, err
		}

		motivation, err := New(rendering.MotivationName, renderingInfo)
		if err != nil {
			return nil, err
		}

		secs := []Entry{
			summary,
			motivation,
		}

		return secs, nil
	*/

	return nil, nil
}

func ForImplementableState(renderingInfo renderingInfoProvider) ([]Entry, error) {

	/*
		// should be read only
		summary, err := New(rendering.SummaryName, renderingInfo)
		if err != nil {
			return nil, err
		}

		// should be read only
		motivation, err := New(rendering.MotivationName, renderingInfo)
		if err != nil {
			return nil, err
		}

		developerGuide, err := New(rendering.DeveloperGuideName, renderingInfo)
		if err != nil {
			return nil, err
		}

		operatorGuide, err := New(rendering.OperatorGuideName, renderingInfo)
		if err != nil {
			return nil, err
		}

		teacherGuide, err := New(rendering.TeacherGuideName, renderingInfo)
		if err != nil {
			return nil, err
		}

		graduationCriteria, err := New(rendering.GraduationCriteriaName, renderingInfo)
		if err != nil {
			return nil, err
		}

		secs := []Entry{
			summary,
			motivation,
			developerGuide,
			operatorGuide,
			teacherGuide,
			graduationCriteria,
		}

		return secs, nil
	*/

	return nil, nil
}
