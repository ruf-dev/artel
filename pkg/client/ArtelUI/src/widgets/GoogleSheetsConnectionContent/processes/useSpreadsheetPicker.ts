import {useRef} from "react"

import {useExternalConnections} from "@/app/hooks/ExternalConnections.ts"
import {useBakeError} from "@/app/hooks/useErrorToast.ts"

type GapiWindow = Window & { gapi: { load: (lib: string, cb: () => void) => void } }

let gapiPromise: Promise<void> | null = null

function loadGapi(): Promise<void> {
    if (gapiPromise) return gapiPromise
    gapiPromise = new Promise<void>((resolve, reject) => {
        const script = document.createElement("script")
        script.src = "https://apis.google.com/js/api.js"
        script.onload = () => (window as unknown as GapiWindow).gapi.load("picker", resolve)
        script.onerror = reject
        document.head.appendChild(script)
    })
    return gapiPromise
}

export function useSpreadsheetPicker() {
    const {addSpreadsheet, getPickerToken} = useExternalConnections()
    const bakeError = useBakeError()
    const pickerOpenRef = useRef(false)

    function openPicker() {
        if (pickerOpenRef.current) return
        pickerOpenRef.current = true

        getPickerToken()
            .then(token => loadGapi().then(() => token))
            .then(token => {
                const picker = new google.picker.PickerBuilder()
                    .addView(new google.picker.DocsView(google.picker.ViewId.SPREADSHEETS))
                    .setOAuthToken(token)
                    .setDeveloperKey(import.meta.env.VITE_GOOGLE_API_KEY ?? "")
                    .setCallback((data: google.picker.ResponseObject) => {
                        const action = data[google.picker.Response.ACTION]
                        if (action === google.picker.Action.PICKED) {
                            const docs = data[google.picker.Response.DOCUMENTS]
                            const file = docs?.[0]
                            if (file) {
                                addSpreadsheet(
                                    file[google.picker.Document.ID],
                                    file[google.picker.Document.NAME] ?? file[google.picker.Document.ID],
                                ).catch(e => bakeError("Failed to add spreadsheet", e))
                            }
                        }
                        if (action === google.picker.Action.PICKED || action === google.picker.Action.CANCEL) {
                            pickerOpenRef.current = false
                        }
                    })
                    .build()
                picker.setVisible(true)
            })
            .catch(e => {
                pickerOpenRef.current = false
                bakeError("Failed to open picker", e)
            })
    }

    return {openPicker}
}
