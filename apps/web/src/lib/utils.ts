import { twMerge } from 'tailwind-merge';

export function cn(...classes: Array<string | false | null | undefined>) {
	return twMerge(classes.filter(Boolean).join(' '));
}

export function capitalize(inputStr: string) {
	if (inputStr && inputStr.length > 0) {
		return inputStr.charAt(0).toUpperCase() + (inputStr.length > 1 ? inputStr.substring(1) : '');
	}

	return inputStr;
}

export function dePluralize(inputStr: string): string {

	const knownEndings = ['s', 'ies'];

	if (inputStr && inputStr.length > 1) {
		const targetKnownEndings = knownEndings
			.filter(ke => inputStr.endsWith(ke))
			.sort((a, b) => b.length - a.length);

		if (targetKnownEndings.length > 0) {
			return inputStr.substring(0, inputStr.length - targetKnownEndings[0].length)
		}

		return inputStr;
	}

	return '';
}
