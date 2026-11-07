import {useEffect, useState} from 'react'
import * as Go from '../wailsjs/go/main/App'
import {BrowserOpenURL} from '../wailsjs/runtime/runtime'

const Icon = ({d}) => (
    <svg className="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"
         strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">{d}</svg>
)
export const UpdateIcon = () => <Icon d={<><path d="M12 3v12"/><path d="m7 10 5 5 5-5"/><path d="M5 21h14"/></>}/>
const XIcon = () => <Icon d={<path d="M18 6 6 18M6 6l12 12"/>}/>

// useUpdate checks GitHub for a newer release once, when the window opens.
export function useUpdate() {
    const [info, setInfo] = useState(null) // {Current, Latest, Newer, Notes, Page}
    useEffect(() => { Go.CheckUpdate().then(setInfo).catch(() => {}) }, [])
    return info
}

// UpdateButton sits in the tab bar; it stands out when a newer version is there.
export function UpdateButton({info, onClick}) {
    const newer = info && info.Newer
    return (
        <button className={'tab help-tab' + (newer ? ' update-ready' : '')} onClick={onClick}
                title={newer ? `Version ${info.Latest} is ready to install` : 'Check for a newer version'}>
            <UpdateIcon/>{newer ? `Update to ${info.Latest}` : 'Updates'}
        </button>
    )
}

// Update checks for the newest release and installs it: the new exe takes
// this one's place, starts, and this window closes. Settings are kept.
export default function Update({onClose}) {
    const [info, setInfo] = useState(null)
    const [error, setError] = useState('')
    const [busy, setBusy] = useState(false)
    const [version, setVersion] = useState('')
    const check = () => {
        setError('')
        setInfo(null)
        Go.CheckUpdate().then(setInfo).catch(e => setError(String(e)))
    }
    useEffect(() => {
        Go.GetVersion().then(setVersion)
        check()
    }, [])
    useEffect(() => {
        const key = e => e.key === 'Escape' && !busy && onClose()
        window.addEventListener('keydown', key)
        return () => window.removeEventListener('keydown', key)
    }, [onClose, busy])

    const install = async () => {
        setBusy(true)
        setError('')
        try {
            await Go.InstallUpdate() // the window closes when it worked
        } catch (e) {
            setError(String(e))
            setBusy(false)
        }
    }

    return (
        <div className="overlay" onClick={() => !busy && onClose()}>
            <div className="guide update" role="dialog" aria-label="Updates" onClick={e => e.stopPropagation()}>
                <div className="guide-head">
                    <UpdateIcon/>
                    <h2>Updates</h2>
                    <span className="spacer"/>
                    <button className="icon-btn" title="Close" aria-label="Close" disabled={busy} onClick={onClose}><XIcon/></button>
                </div>
                <p className="guide-lead">This is Diskless PathFix <b>{version}</b>.</p>

                {error ? <p className="error-text">Could not check for updates: {error}</p>
                    : !info ? <p className="muted">Checking GitHub…</p>
                    : info.Newer ? <>
                        <p>Version <b>{info.Latest}</b> is ready.</p>
                        {info.Notes && <pre className="notes">{info.Notes.trim()}</pre>}
                        <p className="small muted">Update now downloads it, checks it, puts it in place of this exe and
                            opens it. Your packages and settings stay. Run updates on the server (or a super client):
                            on floor PCs the update is gone after a restart.</p>
                    </>
                    : <p className="ok-text">You have the newest version.</p>}

                <div className="guide-foot">
                    {info && info.Page && <button onClick={() => BrowserOpenURL(info.Page)}>Release page</button>}
                    <span className="spacer"/>
                    {busy && <span className="small muted">Downloading…</span>}
                    {info && info.Newer
                        ? <button className="primary" disabled={busy} onClick={install}>Update now</button>
                        : <button disabled={!info && !error} onClick={check}>Check again</button>}
                </div>
            </div>
        </div>
    )
}
