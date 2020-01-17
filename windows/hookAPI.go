// Wrap the windows_hook.go with set of APIs function
package windows

type hook struct {
}

var hookInstance *hook

func Hook() *hook {
	if hookInstance == nil {
		hookInstance = &hook{}
	}

	return hookInstance
}
