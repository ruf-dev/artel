/**
 * Pure function to determine if a swipe gesture should dismiss a toast.
 * @param offsetX The horizontal offset from the drag gesture
 * @param threshold The distance threshold to trigger dismissal (default: 80px)
 * @returns true if the offset meets or exceeds the threshold, false otherwise
 */
export function shouldDismissFromSwipe(offsetX: number, threshold = 80): boolean {
	return Math.abs(offsetX) >= threshold
}
