package models

type Process struct {
	Model

	FileIoInputName  string `json:"file_io_input_name"`
	FileIoOutputName string `json:"file_io_output_name"`

	Definition string `json:"definition"`
}

//publicSpace, err := uuid.NewRandom()
//if err != nil {
//	return nil, errors.Wrap(err, "generate public space failed")
//}
//privateSpace, err := uuid.NewRandom()
//if err != nil {
//	return nil, errors.Wrap(err, "generate private space failed")
//}
//
//if err := s.FileService.CreateFileSpace(publicSpace.String()); err != nil {
//	return nil, err
//}
//if err := s.FileService.CreateFileSpace(privateSpace.String()); err != nil {
//	return nil, err
//}
//
//if p, err = s.SubmissionRepository.CreateProblem(title, locale, publicSpace.String(), privateSpace.String()); err != nil {
//	return nil, err
//}
