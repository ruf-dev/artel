import {create} from "zustand"
import React, {KeyboardEvent} from "react"
import {useCallback} from "react"


export interface DialogManager {
    children: React.JSX.Element[] | null

    IsClickOffClosesDialog: boolean
    closable: boolean

    LockClosing(): void
    UnlockClosing(): void

    SetClosable(v: boolean): void

    OpenDialog(...children: React.JSX.Element[]): void

    CloseDialog(): void
}


export const useDialog =
    create<DialogManager>((set, get) => ({
        children: null,
        IsClickOffClosesDialog: true,
        closable: false,

        LockClosing() {
            set({IsClickOffClosesDialog: false})
        },
        UnlockClosing() {
            set({IsClickOffClosesDialog: true})
        },
        SetClosable(v: boolean) {
            set({closable: v})
        },

        OpenDialog(...children: React.JSX.Element[] ) {
            set({children: children})
        },

        CloseDialog() {
            const {IsClickOffClosesDialog} = get()
            if (IsClickOffClosesDialog) {
                set({children: null, closable: false})
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
