import cls from "@/pages/docs/components/DocsEmptyState/DocsEmptyState.module.css"

// Friendly empty state for `/docs` when the admin hasn't configured a default public
// vault yet — distinct from DocsPage's "This vault could not be found." error, which is
// for a bad/unpublished slug reached via `/docs/:slug`.
export default function DocsEmptyState() {
    return (
        <div className={cls.DocsEmptyStateContainer}>
            No public documentation has been configured yet.
        </div>
    )
}
