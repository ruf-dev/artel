import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import GoogleSheetsIcon from "@/components/ProviderIcon/components/GoogleSheetsIcon"
import TrelloIcon from "@/components/ProviderIcon/components/TrelloIcon"
import MiroIcon from "@/components/ProviderIcon/components/MiroIcon"
import EmailIcon from "@/components/ProviderIcon/components/EmailIcon"
import GitlabIcon from "@/components/ProviderIcon/components/GitlabIcon"
import TelegramIcon from "@/components/ProviderIcon/components/TelegramIcon"
import AnthropicIcon from "@/components/ProviderIcon/components/AnthropicIcon"
import OpenAIIcon from "@/components/ProviderIcon/components/OpenAIIcon"
import OpenRouterIcon from "@/components/ProviderIcon/components/OpenRouterIcon"
import CouchDBIcon from "@/components/ProviderIcon/components/CouchDBIcon"
import S3Icon from "@/components/ProviderIcon/components/S3Icon"
import UnknownProviderIcon from "@/components/ProviderIcon/components/UnknownProviderIcon"

export default function ProviderIcon({provider}: {provider?: ExternalProvider}) {
    switch (provider) {
        case ExternalProvider.EXTERNAL_PROVIDER_GOOGLE_SHEETS:
            return <GoogleSheetsIcon/>
        case ExternalProvider.EXTERNAL_PROVIDER_TRELLO:
            return <TrelloIcon/>
        case ExternalProvider.EXTERNAL_PROVIDER_MIRO:
            return <MiroIcon/>
        case ExternalProvider.EXTERNAL_PROVIDER_EMAIL:
            return <EmailIcon/>
        case ExternalProvider.EXTERNAL_PROVIDER_GITLAB:
            return <GitlabIcon/>
        case ExternalProvider.EXTERNAL_PROVIDER_TELEGRAM:
            return <TelegramIcon/>
        case ExternalProvider.EXTERNAL_PROVIDER_ANTHROPIC:
            return <AnthropicIcon/>
        case ExternalProvider.EXTERNAL_PROVIDER_OPENAI:
            return <OpenAIIcon/>
        case ExternalProvider.EXTERNAL_PROVIDER_OPENROUTER:
            return <OpenRouterIcon/>
        case ExternalProvider.EXTERNAL_PROVIDER_COUCHDB:
            return <CouchDBIcon/>
        case ExternalProvider.EXTERNAL_PROVIDER_S3:
            return <S3Icon/>
        default:
            return <UnknownProviderIcon/>
    }
}
