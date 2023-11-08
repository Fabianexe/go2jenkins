package writer

import (
	"strconv"

	"github.com/Fabianexe/go2jenkins/pkg/entity"
)

func ConvertToCobertura(path string, project *entity.Project) *Coverage {
	pkgs := project.Packages
	coverage := &Coverage{
		Sources: &Sources{
			Sources: []*Source{
				{
					Path: path,
				},
			},
		},
		LineRate:        project.LineCoverage.String(),
		BranchRate:      project.BranchCoverage.String(),
		LinesValid:      project.LineCoverage.ValidString(),
		LinesCovered:    project.LineCoverage.CoveredString(),
		BranchesValid:   project.BranchCoverage.ValidString(),
		BranchesCovered: project.BranchCoverage.CoveredString(),
	}

	packages := &Packages{
		Packages: make([]*Package, 0, len(pkgs)),
	}
	totalComplexity := 0
	for _, pkg := range pkgs {
		packageCov := &Package{
			Name:       pkg.Name,
			LineRate:   pkg.LineCoverage.String(),
			BranchRate: pkg.BranchCoverage.String(),
		}

		packageComplexity := 0

		classes := &Classes{
			Classes: make([]*Class, 0, len(pkg.Files)),
		}

		for _, file := range pkg.Files {
			class := &Class{
				Name:       file.Name,
				Filename:   file.FilePath,
				LineRate:   file.LineCoverage.String(),
				BranchRate: file.BranchCoverage.String(),
			}

			classComplexity := 0

			methods := &Methods{
				Methods: make([]*Method, 0, len(file.Methods)),
			}

			classLines := &Lines{
				Lines: make([]*Line, 0, 1024),
			}

			for _, method := range file.Methods {
				xmlMethod := &Method{
					Name:       method.Name,
					LineRate:   method.LineCoverage.String(),
					BranchRate: method.BranchCoverage.String(),
					Complexity: strconv.Itoa(method.Complexity),
				}

				totalComplexity += method.Complexity
				packageComplexity += method.Complexity
				classComplexity += method.Complexity

				methodsLines := &Lines{
					Lines: make([]*Line, 0, len(method.Lines)),
				}

				branchstarts := make([]int, 0, len(method.Branches))
				for _, branch := range method.Branches[1:] {
					branchstarts = append(branchstarts, branch.StartLine)
				}

				for _, line := range method.Lines {
					xmlLine := &Line{
						Number: strconv.Itoa(line.Number),
						Hits:   strconv.Itoa(line.CoverageCount),
						Branch: "false",
					}
					for i := 0; i < len(branchstarts); {
						if line.Number >= branchstarts[i] {
							xmlLine.Branch = "true"
							branchstarts = append(branchstarts[:i], branchstarts[i+1:]...)
							continue
						}

						i++
					}
					methodsLines.Lines = append(methodsLines.Lines, xmlLine)
					classLines.Lines = append(classLines.Lines, xmlLine)
				}

				if len(methodsLines.Lines) != 0 {
					xmlMethod.Lines = methodsLines
				}

				methods.Methods = append(methods.Methods, xmlMethod)
			}
			if len(methods.Methods) != 0 {
				class.Methods = methods
			}

			if len(classLines.Lines) != 0 {
				class.Lines = classLines
			}

			class.Complexity = strconv.Itoa(classComplexity)

			classes.Classes = append(classes.Classes, class)
		}

		if len(classes.Classes) != 0 {
			packageCov.Classes = classes
		}

		packageCov.Complexity = strconv.Itoa(packageComplexity)

		packages.Packages = append(packages.Packages, packageCov)
	}
	if len(packages.Packages) != 0 {
		coverage.Packages = packages
	}

	coverage.Complexity = strconv.Itoa(totalComplexity)

	return coverage
}
