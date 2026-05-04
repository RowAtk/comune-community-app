export const communitySections = [
	{ href: '', label: 'Overview', description: 'Snapshot of the active community' },
	{ href: 'units', label: 'Units', description: 'Physical spaces, blocks, and occupancy anchors' },
	{
		href: 'households',
		label: 'Households',
		description: 'Family and occupancy groups mapped to units'
	},
	{
		href: 'residents',
		label: 'Residents',
		description: 'Owners, tenants, dependents, and contacts'
	},
	{
		href: 'invoicing',
		label: 'Invoicing',
		description: 'Maintenance billing, overdue households, and resident due-date visibility'
	}
] as const;
