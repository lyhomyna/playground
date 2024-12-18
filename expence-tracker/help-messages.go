package main

func helpText() string {
    return ` 
    Usage: expence-tracker <command>
    Where <command> can be add, list, summary, delete.
    `
}

func helpAddText() string {
    return `
Usage example: 
    $ expence-tracer add --description "Lunch" --amount 20
`
}

func helpDeleteText() string {
    return `
Usage example: 
    $ expence-tracer delete --id 3
`
}

func helpSummaryText() string {
    return `
Usage example: 
    $ expence-tracer summary		// to print all summary
    $ expence-tracer summary --month 8  // to print summary by month
`
}
