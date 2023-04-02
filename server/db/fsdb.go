package db

func SavePipeline() {

}

type ProjectInfo struct {
	Name      string
	Pipelines []string
}

func GetProjectInfo() []ProjectInfo {
	return MockProjectInfo()
}

func MockProjectInfo() []ProjectInfo {
	p1 := ProjectInfo{
		Name:      "project-1",
		Pipelines: []string{"pipeline-1", "pipeline-2"},
	}
	p2 := ProjectInfo{
		Name:      "project-2",
		Pipelines: []string{"pipeline-3"},
	}
	return []ProjectInfo{p1, p2}
}
