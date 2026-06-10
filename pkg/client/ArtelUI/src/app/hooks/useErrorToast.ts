import {useToaster} from "@vervstack/chures"

export function useBakeError() {
    const {bake} = useToaster()
    return function(title: string, err: unknown) {
        bake({
            title,
            description: err instanceof Error ? err.message : "Unexpected error",
            level: "Error",
        })
    }
}
