export interface Update {
	Type: string
}

export interface AddProcess extends Update {
	CBT: number
	Name: string
	AT: number
	QI: number
}

export interface CPUUpdate extends Update {
	Name: string
	Start: number
	End: number
	QI: number
}