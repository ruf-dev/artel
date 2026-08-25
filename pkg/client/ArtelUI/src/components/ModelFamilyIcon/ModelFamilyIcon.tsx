import AnthropicIcon from "@/components/ProviderIcon/components/AnthropicIcon"
import DeepSeekIcon from "@/components/ProviderIcon/components/DeepSeekIcon"
import MetaLlamaIcon from "@/components/ProviderIcon/components/MetaLlamaIcon"
import GoogleIcon from "@/components/ProviderIcon/components/GoogleIcon"
import MistralIcon from "@/components/ProviderIcon/components/MistralIcon"
import QwenIcon from "@/components/ProviderIcon/components/QwenIcon"
import XAIIcon from "@/components/ProviderIcon/components/XAIIcon"
import UnknownProviderIcon from "@/components/ProviderIcon/components/UnknownProviderIcon"

export default function ModelFamilyIcon({family}: {family: string}) {
    const normalizedFamily = family.toLowerCase()

    switch (normalizedFamily) {
        case "anthropic":
            return <AnthropicIcon/>
        case "deepseek":
            return <DeepSeekIcon/>
        case "meta-llama":
            return <MetaLlamaIcon/>
        case "google":
            return <GoogleIcon/>
        case "mistralai":
            return <MistralIcon/>
        case "qwen":
            return <QwenIcon/>
        case "x-ai":
            return <XAIIcon/>
        default:
            return <UnknownProviderIcon/>
    }
}
