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
    $ expence-tracker add --description "Lunch" --amount 20
`
}

func helpDeleteText() string {
    return `
Usage example: 
    $ expence-tracker delete        // delete all
    $ expence-tracker delete --id 3 // deleting by id
`
}

func helpSummaryText() string {
    return `
Usage example: 
    $ expence-tracker summary		 // print all summary
    $ expence-tracker summary --month 8  // print summary by month
`
}
