import {useCallback} from "react"
import {useToaster} from "@vervstack/chures"
import {grpcErrorMessage, isGrpcError} from "@/processes/grpcErrors"

export function useBakeError() {
    const {bake} = useToaster()
    return useCallback((title: string, err: unknown) => {
        let description: string
        if (isGrpcError(err)) {
            description = grpcErrorMessage(err)
        } else if (err instanceof Error) {
            description = err.message
        } else {
            description = "Unexpected error"
        }
        bake({
            title,
            description,
            level: "Error",
        })
    }, [bake])
}
