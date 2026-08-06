import {ExternalProvider} from "@/app/api/artel/external_connections.pb.ts"
import GoogleSheetsIcon from "@/components/ProviderIcon/components/GoogleSheetsIcon"
import TrelloIcon from "@/components/ProviderIcon/components/TrelloIcon"
import MiroIcon from "@/components/ProviderIcon/components/MiroIcon"
import EmailIcon from "@/components/ProviderIcon/components/EmailIcon"
import GitlabIcon from "@/components/ProviderIcon/components/GitlabIcon"
import TelegramIcon from "@/components/ProviderIcon/components/TelegramIcon"
import AnthropicIcon from "@/components/ProviderIcon/components/AnthropicIcon"
import OpenAIIcon from "@/components/ProviderIcon/components/OpenAIIcon"
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
        default:
            return <UnknownProviderIcon/>
    }
}
