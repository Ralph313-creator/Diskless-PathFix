import {useEffect, useState} from 'react'
import * as Go from '../wailsjs/go/main/App'
import {ClipboardSetText} from '../wailsjs/runtime/runtime'

const Icon = ({d}) => (
    <svg className="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"
         strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">{d}</svg>
)
export const FileIcon = () => <Icon d={<><path d="M14 3H7a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V8z"/><path d="M14 3v5h5M9 13h6M9 17h6"/></>}/>
const XIcon = () => <Icon d={<path d="M18 6 6 18M6 6l12 12"/>}/>

// copyText puts text on the clipboard and says whether it worked.
export async function copyText(text) {
    try {
        if (await ClipboardSetText(text)) return true
    } catch { /* try the browser's clipboard */ }
    try {
        await navigator.clipboard.writeText(text)
        return true
    } catch {
        return false
    }
}

// PublicTxt shows every line for CCBoot's C:\CCBoot\Cmd\public.txt (each
// portable app, then GameFix.bat), to copy and paste there. The app never
// writes that file itself: Open public.txt only opens it in Notepad.
export default function PublicTxt({onClose}) {
    const [text, setText] = useState(null) // null while loading
    const [msg, setMsg] = useState(null)   // {ok, text}
    useEffect(() => { Go.PublicTxt().then(setText) }, [])
    useEffect(() => {
        const key = e => e.key === 'Escape' && onClose()
        window.addEventListener('keydown', key)
        return () => window.removeEventListener('keydown', key)
    }, [onClose])

    const lines = (text || '').split('\n').filter(l => l.trim()).length
    const copy = async () => setMsg(await copyText(text)
        ? {ok: true, text: `Copied ${lines} line${lines === 1 ? '' : 's'}. Now paste them into public.txt.`}
        : {ok: false, text: 'Could not copy. Select the lines and press Ctrl+C.'})
    const open = async () => {
        try {
            await Go.OpenPublicTxt()
        } catch (e) {
            setMsg({ok: false, text: String(e)})
        }
    }

    return (
        <div className="overlay" onClick={onClose}>
            <div className="guide" role="dialog" aria-label="Lines for public.txt" onClick={e => e.stopPropagation()}>
                <div className="guide-head">
                    <FileIcon/>
                    <h2>Lines for public.txt</h2>
                    <span className="spacer"/>
                    <button className="icon-btn" title="Close" aria-label="Close" onClick={onClose}><XIcon/></button>
                </div>
                <p className="guide-lead">Floor PCs only set up what is in CCBoot's <code>C:\CCBoot\Cmd\public.txt</code>.
                    These are the lines for every portable app and SymbolicLink in Diskless PathFix.</p>

                {text === null ? <p className="muted">Loading…</p>
                    : text === '' ? (
                        <div className="empty">
                            <h2>Nothing to paste yet</h2>
                            <p className="muted">Add a portable app, or set up SymbolicLink, first. Their lines show up here.</p>
                        </div>
                    ) : (
                        <textarea className="lines" readOnly value={text.trimEnd()} rows={Math.min(Math.max(lines, 2), 8)}
                                  spellCheck={false} onFocus={e => e.target.select()}/>
                    )}

                <ol className="how">
                    <li>Click <b>Copy all</b>.</li>
                    <li>Click <b>Open public.txt</b>. If it already has RapTech lines (with <code>\Portable\</code> or
                        <code> GameFix.bat</code>), delete them, so each line is there only once. Then paste.</li>
                    <li>Save the file (Ctrl+S), and restart one floor PC to test.</li>
                </ol>
                <p className="small muted">Apps run with <b>-startup</b>, when a customer logs on. GameFix.bat runs
                    with <b>-service</b>: as SYSTEM, with admin rights, before anyone logs on. The lines use this PC's
                    drive letters; if the game disk has another letter on the floor PCs, change it after pasting.
                    Don't add other text to the lines.</p>

                <div className="guide-foot">
                    {msg && <span className={'small ' + (msg.ok ? 'ok-text' : 'error-text')}>{msg.text}</span>}
                    <span className="spacer"/>
                    <button onClick={open}><FileIcon/>Open public.txt</button>
                    <button className="primary" disabled={!text} onClick={copy}>Copy all</button>
                </div>
            </div>
        </div>
    )
}
