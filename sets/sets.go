package sets

func LoadModule(modulename string) *[]JSONCommand {
	switch modulename {
	case "general":
		return newGenModule()
	case "moderation":
		return newModModule()
	case "fun":
		return newFunModule()
	}
}
