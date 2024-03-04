package internal

// AllTiers is a list of all tiers that can be selected.
// For descriptions, please see: docs/tiers/README.md
var AllTiers = []string{
	"Undefined",
	"Tier 1 - Critical",
	"Tier 2 - Important",
	"Tier 3 - Ancillary",
	"Tier 4 - Internal",
}

// AllCapabilities is a list of all capabilities that can be selected.
// For descriptions, please see: docs/capabilities/README.md
var AllCapabilities = []string{
	"Undefined",
	"Order Injection",
	"Menu Sync from POS",
	"Menu Sync to JE",
	"Menu Sync to TA",
	"GFO Inventory",
	"Menu Curation",
	"PLU Lookup",
	"Basket Calculation",
	"Item Availability",
	"Delivery",
	"Event Bridge",
	"Restaurants",
}

var AllTeams = []string{
	"fd-new-domains",
	"core-eng-rel",
	"jetc-operations",
	"jetc-groceries",
	"jetc-tooling",
	"jetc-canada",
}
