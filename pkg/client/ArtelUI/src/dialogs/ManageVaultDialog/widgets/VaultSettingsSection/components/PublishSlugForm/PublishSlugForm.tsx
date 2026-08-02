import {useState} from "react"
import {Button} from "@vervstack/chures"

import Input from "@/components/atoms/Input/Input.tsx"
import cls from
    "@/dialogs/ManageVaultDialog/widgets/VaultSettingsSection/components/PublishSlugForm/PublishSlugForm.module.css"

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
    vaultName: string
    submitting: boolean
    onPublish: (slug: string) => void
    onCancel: () => void
}

export default function PublishSlugForm({vaultName, submitting, onPublish, onCancel}: Props) {
    const [slug, setSlug] = useState(() => slugify(vaultName))

    const error = validateSlug(slug)

    return (
        <div className={cls.PublishSlugFormContainer}>
            <label className={cls.Field}>
                <span className={cls.FieldLabel}>Public URL slug</span>
                <Input
                    placeholder="my-vault"
                    value={slug}
                    setValue={setSlug}
                    error={error}
                    disabled={submitting}
                    autoComplete="off"
                />
            </label>
            <div className={cls.Actions}>
                <Button variant="ghost" onClick={onCancel} disabled={submitting}>
                    Cancel
                </Button>
                <Button variant="primary" onClick={() => onPublish(slug)} disabled={submitting || !!error}>
                    {submitting ? "Publishing…" : "Publish"}
                </Button>
            </div>
        </div>
    )
}
