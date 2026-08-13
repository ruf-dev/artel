import {useState} from "react"
import {Button} from "@vervstack/chures"

import {VaultItem} from "@/app/api/artel/vaults.pb.ts"
import {useVaultMutations} from "@/app/hooks/Vaults.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"
import Input from "@/components/atoms/Input/Input.tsx"
import cls from "@/pages/admin/components/SettingsTab/components/DefaultDocsVaultSection/components/PrivateVaultRow/PrivateVaultRow.module.css" // eslint-disable-line max-len

// Mirrors dialogs/ManageVaultDialog/widgets/VaultSettingsSection/components/PublishSlugForm's
// slug validation inline rather than importing it — that component is local to its own
// dialog subtree by convention, and this duplication is small enough not to warrant
// promoting a shared abstraction for it.
const SLUG_PATTERN = /^[a-z0-9]+(-[a-z0-9]+)*$/

function slugify(name: string): string {
    return name
        .toLowerCase()
        .replace(/[^a-z0-9]+/g, "-")
        .replace(/^-+|-+$/g, "")
}

function validateSlug(slug: string): string | undefined {
    if (!slug) {
        return "Slug is required"
    }
    if (!SLUG_PATTERN.test(slug)) {
        return "Use lowercase letters, numbers, and single hyphens between words, e.g. my-vault"
    }
    return undefined
}

interface Props {
    vault: VaultItem
    onPublished: (vault: VaultItem) => void
}

export default function PrivateVaultRow({vault, onPublished}: Props) {
    const {publish} = useVaultMutations()
    const bakeError = useBakeError()
    const [formOpen, setFormOpen] = useState(false)
    const [slug, setSlug] = useState(() => slugify(vault.name ?? ""))
    const [publishing, setPublishing] = useState(false)

    const error = validateSlug(slug)

    function handlePublish() {
        setPublishing(true)
        publish(vault.id ?? "", slug)
            .then(published => {
                onPublished(published)
                setFormOpen(false)
            })
            .catch(err => bakeError("Failed to publish vault", err))
            .finally(() => setPublishing(false))
    }

    return (
        <div className={cls.PrivateVaultRowContainer}>
            <div className={cls.Row}>
                <span className={cls.Name}>{vault.name}</span>
                <Button variant="secondary" onClick={() => setFormOpen(open => !open)} disabled={publishing}>
                    {formOpen ? "Cancel" : "Publish"}
                </Button>
            </div>
            {formOpen && (
                <div className={cls.SlugForm}>
                    <Input
                        placeholder="my-vault"
                        value={slug}
                        setValue={setSlug}
                        error={error}
                        disabled={publishing}
                        autoComplete="off"
                    />
                    <Button variant="primary" onClick={handlePublish} disabled={publishing || !!error}>
                        {publishing ? "Publishing…" : "Confirm"}
                    </Button>
                </div>
            )}
        </div>
    )
}
