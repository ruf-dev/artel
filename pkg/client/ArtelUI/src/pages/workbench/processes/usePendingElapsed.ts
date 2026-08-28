import {useEffect, useState} from "react"

export const SLOW_AFTER_MS = 20000
export const STUCK_AFTER_MS = 75000

export type PendingBucket = "normal" | "slow" | "stuck"

// Tracks elapsed time while `active` is true. Returns "normal" for the first 20s,
// "slow" from 20s to 75s, and "stuck" after 75s. Resets the timer whenever `active`
// becomes true or whenever `resetKey` changes while active.
export function usePendingElapsed(active: boolean, resetKey: unknown): PendingBucket {
    const [bucket, setBucket] = useState<PendingBucket>("normal")

    useEffect(() => {
        if (!active) {
            setBucket("normal")
            return
        }

        // Start at "normal" when active starts
        setBucket("normal")

        // Set the "slow" timer
        const slowTimer = setTimeout(() => {
            setBucket("slow")
        }, SLOW_AFTER_MS)

        // Set the "stuck" timer
        const stuckTimer = setTimeout(() => {
            setBucket("stuck")
        }, STUCK_AFTER_MS)

        return () => {
            clearTimeout(slowTimer)
            clearTimeout(stuckTimer)
        }
    }, [active, resetKey])

    return bucket
}
