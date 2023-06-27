export interface Update {
	Type: string
}

export interface AddProcess extends Update {
	CBT: number
	Name: string
	AT: number
	AtStr: string | undefined
	QI: number
}

export interface CPUUpdate extends Update {
	Name: string
	Start: number
	StartStr: string | undefined
	End: number
	EndStr: string | undefined
	QI: number
}