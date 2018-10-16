package section

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/calebamiles/keps/pkg/keps/section/internal/rendering"
)

// New creates a new rendered section, either by reading existing content on disk
// or by rendering a new section of the given sectionName. Sections which are always
// regenerated (e.g. README.md, _index.json) should not use new
func New(sectionName string, info renderingInfoProvider) (section, error) {
	var sec section
	var creator *creator

	// don't render sections which are auto generated and require other section information
	switch sectionName {
	case rendering.ReadmeName:
		return nil, errors.New("cannot render README section using section.New(), use section.NewReadme()")
	default:
		creator = creatorFor[sectionName]
		if creator == nil {
			return nil, fmt.Errorf("no top level KEP section: %s exists", sectionName)
		}

	}

	contentDir := info.ContentDir()
	loc := filepath.Join(contentDir, creator.filename)
	readBytes, err := ioutil.ReadFile(loc)

	switch {
	case os.IsNotExist(err):
		contentBytes, contentErr := creator.render(info)
		if contentErr != nil {
			return nil, contentErr
		}

		sec = &persistableSection{
			commonSectionInfo: &commonSectionInfo{
				filename:   creator.filename,
				name:       sectionName,
				content:    contentBytes,
				contentDir: contentDir,
			},
		}

	default:
		sec = &readOnlySection{
			commonSectionInfo: &commonSectionInfo{
				filename:   creator.filename,
				name:       sectionName,
				content:    readBytes,
				contentDir: contentDir,
			},
		}
	}

	return sec, nil
}

func NewReadme(info rendering.InfoAndSectionProvider) (section, error) {
	readmeBytes, err := rendering.NewReadme(info)
	if err != nil {
		return nil, err
	}

	sec := &persistableSection{
		commonSectionInfo: &commonSectionInfo{
			filename:   rendering.ReadmeFilename,
			name:       rendering.ReadmeName,
			contentDir: info.ContentDir(),
			content:    readmeBytes,
		},
	}

	return sec, nil
}

type renderer func(rendering.InfoProvider) ([]byte, error)

type creator struct {
	render   renderer
	filename string
}

var creatorFor = map[string]*creator{
	rendering.SummaryName: &creator{
		render:   rendering.NewSummary,
		filename: rendering.SummaryFilename,
	},

	rendering.MotivationName: &creator{
		render:   rendering.NewMotivation,
		filename: rendering.MotivationFilename,
	},

	rendering.DeveloperGuideName: &creator{
		render:   rendering.NewDeveloperGuide,
		filename: rendering.DeveloperGuideFilename,
	},

	rendering.OperatorGuideName: &creator{
		render:   rendering.NewOperatorGuide,
		filename: rendering.OperatorGuideFilename,
	},

	rendering.TeacherGuideName: &creator{
		render:   rendering.NewTeacherGuide,
		filename: rendering.TeacherGuideFilename,
	},

	rendering.GraduationCriteriaName: &creator{
		render:   rendering.NewGraduationCriteria,
		filename: rendering.GraduationCriteriaFilename,
	},
}
