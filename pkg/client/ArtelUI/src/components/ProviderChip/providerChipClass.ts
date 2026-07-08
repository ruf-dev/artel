import cls from "@/components/ProviderChip/ProviderChip.module.css"

import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"

export const PROVIDER_CHIP_CLASS: Partial<Record<ExternalProvider, string>> = {
    [ExternalProvider.EXTERNAL_PROVIDER_GOOGLE_SHEETS]: cls.SheetsChip,
    [ExternalProvider.EXTERNAL_PROVIDER_TRELLO]: cls.TrelloChip,
    [ExternalProvider.EXTERNAL_PROVIDER_MIRO]: cls.MiroChip,
    [ExternalProvider.EXTERNAL_PROVIDER_GITLAB]: cls.GitlabChip,
}
