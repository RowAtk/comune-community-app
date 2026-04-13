import type { CommunityResource } from '$lib/api/types';
import type { Snippet } from 'svelte';

export type ResourceCollectionFormState<TValues extends Record<string, unknown> = Record<string, unknown>> = {
	createError?: string;
	createSuccess?: string;
	values?: Partial<TValues>;
};

export type ResourceCollectionProps<
	T extends CommunityResource,
	TValues extends Record<string, unknown> = Record<string, unknown>
> = {
	resourceName: string;
	resourceSingleName?: string;
	resourceHeading: string;
	resourceDescription: string;
	creationDescription: string;
	resourceList: T[];
	createForm?: ResourceCollectionFormState<TValues> | null;
	createFormBody?: Snippet;
	createAction?: string;
	apiError?: string | null;
	resourceEmptyCreateDescription?: string;
	headerActions?: Snippet;
	groupBy?: (resource: T) => string | null | undefined;
	listItemCard: Snippet<[T]>;
};
