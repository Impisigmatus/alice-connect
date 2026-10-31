package usecase

func (uc *UseCase) SerialCommand(cmd string) (string, error) {
	return uc.serialDriverAPI.PostSerialCommand(cmd)
}
