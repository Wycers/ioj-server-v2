package manager

type Handler func(element *ProcessElement) (bool, error)
