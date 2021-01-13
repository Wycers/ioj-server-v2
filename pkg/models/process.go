package models

type Process struct {
	Model

	FileIoInputName  string `json:"file_io_input_name"`
	FileIoOutputName string `json:"file_io_output_name"`

	Definition string `json:"definition"`
}
