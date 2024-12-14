package allocator

var (
	RegistersX86 = [16]string{
		"r15", "r14", "r13", "r12", "r11", "r10",
		"r9", "r8", "rsp", "rbp", "rdi", "rsi",
		"rdx", "rcx", "rbx", "rax"}
	contentTypeName = map[RegisterContentType]string{
		NumberRegister:   "number",
		VariableRegister: "variable",
	}
)

type RegisterContentType int

type Register struct {
	isFree      bool
	name        string
	content     string
	contentType RegisterContentType
	isLoaded    bool
}

const (
	NumberRegister = iota
	VariableRegister
)

func (contentType RegisterContentType) String() {
	return
}

func (register *Register) setRegisterStatus(status bool) {
	register.isFree = status
}

func (register *Register) SetIsLoaded(val bool) {
	register.isLoaded = val
}

func (register *Register) SetType(val RegisterContentType) {
	register.contentType = val
}

func (register *Register) GetIsLoaded() bool {
	return register.isLoaded
}

func (register *Register) GetContent() string {
	return register.content
}

func (register *Register) GetType() RegisterContentType {
	return register.contentType
}

func (register *Register) GetName() string {
	return register.name
}

func (register *Register) Load(content string) {
	register.content = content
}
