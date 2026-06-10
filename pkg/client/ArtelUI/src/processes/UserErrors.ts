import {getPreconditionViolations, GrpcStatusError} from "@/processes/grpcErrors.ts";
import {UserErrors} from "@/app/api/artel";

export function isMissingSubscription(err: GrpcStatusError): boolean {
    return getPreconditionViolations(err).
    find(v => v.description == UserErrors.NoSubscription) !== undefined
}
