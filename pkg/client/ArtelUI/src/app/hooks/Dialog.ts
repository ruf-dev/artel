import {create} from "zustand"
import React, {KeyboardEvent} from "react"
import {useCallback} from "react"


export interface DialogManager {
    children: React.JSX.Element[] | null

    IsClickOffClosesDialog: boolean

    LockClosing(): void

    UnlockClosing(): void

    OpenDialog(...children: React.JSX.Element[]): void

    CloseDialog(): void
}


export const useDialog =
    create<DialogManager>((set, get) => ({
        children: null,
        IsClickOffClosesDialog: true,

        LockClosing() {
            set({IsClickOffClosesDialog: false})
        },
        UnlockClosing() {
            set({IsClickOffClosesDialog: true})
        },

        OpenDialog(...children: React.JSX.Element[] ) {
            set({children: children})
        },

        CloseDialog() {
            const {IsClickOffClosesDialog} = get()
            if (IsClickOffClosesDialog) {
                set({children: null})
            }
        }
    }))

export function useDialogKeyboard(onEnter: () => void) {
    return useCallback((e: KeyboardEvent<HTMLInputElement>) => {
        if (e.key === "Enter") void onEnter()
        if (e.key === "Escape") {
            const {CloseDialog} = useDialog.getState()
            CloseDialog()
        }
    }, [onEnter])
}
