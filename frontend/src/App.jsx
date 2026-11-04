import {useEffect, useRef, useState} from 'react'
import * as Go from '../wailsjs/go/main/App'
import {EventsOn} from '../wailsjs/runtime/runtime'

const Icon = ({d}) => (
    <svg className="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"
         strokeLinecap="round" strokeLinejoin="round">{d}</svg>
)
const FolderIcon = () => <Icon d={<path d="M3 7a2 2 0 0 1 2-2h4l2 2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>}/>
const KeyIcon = () => <Icon d={<><circle cx="8" cy="15" r="4"/><path d="m11 12 9-9M17 6l3 3"/></>}/>
const RefreshIcon = () => <Icon d={<><path d="M21 12a9 9 0 1 1-3-6.7L21 8"/><path d="M21 3v5h-5"/></>}/>
const XIcon = () => <Icon d={<path d="M18 6 6 18M6 6l12 12"/>}/>
const SearchIcon = () => <Icon d={<><circle cx="11" cy="11" r="7"/><path d="m21 21-4.3-4.3"/></>}/>
const LinkIcon = () => <Icon d={<><path d="M10 13a5 5 0 0 0 7 0l3-3a5 5 0 0 0-7-7l-1 1"/><path d="M14 11a5 5 0 0 0-7 0l-3 3a5 5 0 0 0 7 7l1-1"/></>}/>

const SCHEDULES = [[0, 'Off'], [15, 'Every 15 minutes'], [30, 'Every 30 minutes'], [60, 'Every hour'], [360, 'Every 6 hours'], [1440 - 1, 'Once a day']]

const stateText = {0: 'not linked', 1: 'linked', 2: 'links elsewhere', 3: 'real folder (not linked)', 4: 'a file is in the way'}

export default function App() {
    const [root, setRoot] = useState('')
    const [packages, setPackages] = useState([])
    const [selected, setSelected] = useState('')
    const [busy, setBusy] = useState('')
    const [toast, setToast] = useState(null)
    const [creating, setCreating] = useState('')
    const [hunt, setHunt] = useState(null) // {dir, state}
    const huntDir = useRef('')

    useEffect(() => EventsOn('hunt', st => setHunt({dir: huntDir.current, state: st})), [])

    const refresh = async () => {
        const list = (await Go.List()) || []
        setPackages(list)
        if (!list.find(p => p.Dir === selected) && list.length) setSelected(list[0].Dir)
    }
    useEffect(() => { Go.GetRoot().then(r => { setRoot(r); refresh() }) }, [])

    // run does one action with a busy label and shows how it went.
    const run = async (label, fn, okText) => {
        setBusy(label)
        try {
            const out = await fn()
            setToast({ok: true, text: okText || (typeof out === 'string' && out ? lastLine(out) : 'Done')})
        } catch (e) {
            setToast({ok: false, text: String(e)})
        } finally {
            setBusy('')
            refresh()
        }
    }

    const pkg = packages.find(p => p.Dir === selected)

    return (
        <div className="app">
            <div className="topbar">
                <div className="brand"><LinkIcon/><h1>RapTech Portable</h1></div>
                <div className="main-folder">
                    <span>Apps folder</span>
                    <span className="path">{root || 'not chosen'}</span>
                    <button onClick={async () => { setRoot(await Go.ChooseRoot()); refresh() }}>Choose…</button>
                </div>
            </div>
            <div className="body">
                <div className="sidebar">
                    <div className="sidebar-title">Portable apps</div>
                    <div className="game-list">
                        {packages.map(p => (
                            <button key={p.Dir} className={'game-item' + (p.Dir === selected ? ' active' : '')}
                                    onClick={() => { setSelected(p.Dir); setCreating('') }}>
                                <span className="game-name">{p.Name}</span>
                                <span className="count">{p.Error ? 'error' : p.Version ? p.Version.replace('version-', '').slice(0, 6) : ''}</span>
                            </button>
                        ))}
                    </div>
                    <div className="new-game" style={{flexDirection: 'column'}}>
                        <button className="primary" disabled={!!busy || !root} onClick={() => setCreating('roblox')}>Add Roblox</button>
                        <button disabled={!!busy || !root} onClick={() => setCreating('app')}>Add another app</button>
                    </div>
                </div>
                <div className="content">
                    {creating === 'roblox' && <NewRoblox busy={busy} run={run} done={(dir) => { setCreating(''); setSelected(dir) }} root={root}/>}
                    {creating === 'app' && <NewApp busy={busy} run={run} done={(dir) => { setCreating(''); setSelected(dir) }} root={root}/>}
                    {!creating && !pkg && (
                        <div className="empty">
                            <h2>No portable apps yet</h2>
                            <p className="muted">Choose the apps folder on the game disk, such as G:\Portable, then add Roblox or another app.</p>
                        </div>
                    )}
                    {!creating && pkg && <Package pkg={pkg} busy={busy} run={run}
                        hunt={hunt && hunt.dir === pkg.Dir ? hunt.state : null}
                        startHunt={async () => {
                            huntDir.current = pkg.Dir
                            try {
                                const exe = await Go.StartHunt(pkg.Dir)
                                if (exe) setHunt({dir: pkg.Dir, state: {exe, running: true, programs: [], folders: [], live: 0}})
                            } catch (e) { setToast({ok: false, text: String(e)}) }
                        }}
                        closeHunt={() => setHunt(null)}/>}
                </div>
            </div>
            <div className="footer">
                <span>{busy ? busy + '…' : 'Updates only stick when they run here on the server: floor PCs throw away game disk writes at restart.'}</span>
            </div>
            {toast && (
                <div className={'toast ' + (toast.ok ? 'ok' : 'error')}>
                    <span>{toast.text}</span>
                    <button className="icon-btn" onClick={() => setToast(null)}><XIcon/></button>
                </div>
            )}
        </div>
    )
}

function lastLine(s) {
    const lines = s.trim().split('\n')
    return lines[lines.length - 1].replace(/^\S+ \S+ \S+ \[\d+\] /, '')
}

function NewRoblox({busy, run, done, root}) {
    const [perUser, setPerUser] = useState(false)
    useEffect(() => { Go.RobloxPerUser().then(setPerUser) }, [])
    return (
        <div className="card">
            <div className="card-head"><h2>Add Roblox</h2></div>
            <p className="muted">Makes <code>{root}\Roblox</code> and downloads the newest Roblox Player into it straight from Roblox.
                Customers' Roblox logins are emptied when Roblox closes.</p>
            <label className="watch"><input type="radio" checked={!perUser} onChange={() => setPerUser(false)}/>
                Installed for everyone (C:\Program Files\Roblox) — like this server</label>
            <label className="watch"><input type="radio" checked={perUser} onChange={() => setPerUser(true)}/>
                Installed for one user (%LOCALAPPDATA%\Roblox)</label>
            <p className="small muted">Pick how Roblox is installed in the client image, so its folder gets linked.</p>
            <div className="hunt-actions">
                <button className="primary" disabled={!!busy}
                        onClick={() => run('Downloading Roblox', () => Go.MakeRoblox(perUser)).then(() => done(root + '\\Roblox'))}>
                    Make and download
                </button>
            </div>
        </div>
    )
}

function NewApp({busy, run, done, root}) {
    const [name, setName] = useState('')
    const [appDir, setAppDir] = useState('')
    const [exe, setExe] = useState('')
    const pickExe = async () => {
        const f = await Go.PickExe()
        if (!f) return
        setExe(f)
        if (!appDir) setAppDir(f.substring(0, f.lastIndexOf('\\')))
        if (!name) setName(f.substring(f.lastIndexOf('\\') + 1).replace(/\.exe$/i, ''))
    }
    return (
        <div className="card">
            <div className="card-head"><h2>Add another app</h2></div>
            <p className="muted">The app's folder is copied into <code>{root}\{name || '<name>'}</code>. Then add the folders
                and registry keys it keeps outside its own folder.</p>
            <div className="field"><span className="field-label">Program</span>
                <input value={exe} onChange={e => setExe(e.target.value)} placeholder="C:\Games\Game\game.exe"/>
                <button onClick={pickExe}>Browse…</button></div>
            <div className="field"><span className="field-label">Folder</span>
                <input value={appDir} onChange={e => setAppDir(e.target.value)} placeholder="the app's whole folder"/>
                <button onClick={async () => { const d = await Go.PickFolder("The app's folder"); if (d) setAppDir(d) }}>Browse…</button></div>
            <div className="field"><span className="field-label">Name</span>
                <input value={name} onChange={e => setName(e.target.value)} placeholder="Game"/></div>
            <div className="hunt-actions">
                <button className="primary" disabled={!!busy || !name}
                        onClick={() => run('Copying ' + name, () => Go.CreatePackage(name, appDir, exe), 'Made ' + name).then(() => done(root + '\\' + name))}>
                    Make package
                </button>
            </div>
        </div>
    )
}

function Package({pkg, busy, run, hunt, startHunt, closeHunt}) {
    const [link, setLink] = useState('')
    const [session, setSession] = useState('')
    const [key, setKey] = useState('')
    if (pkg.Error) {
        return <div className="card"><h2>{pkg.Name}</h2><p className="error-text">{pkg.Error}</p></div>
    }
    return (
        <>
            <div className="game-head">
                <div>
                    <h2>{pkg.Name}</h2>
                    <p className="small muted">{pkg.Dir}</p>
                </div>
                <div className="actions">
                    <button onClick={() => Go.Open(pkg.Dir)}><FolderIcon/>Open</button>
                </div>
            </div>

            {pkg.Update && (
                <div className="card">
                    <div className="card-head"><RefreshIcon/><span className="name">Updates</span><span className="spacer"/>
                        {pkg.Version && <span className="pill ok">{pkg.Version}</span>}</div>
                    <p className="small muted">{pkg.Update.kind === 'roblox' ? 'Downloads each new Roblox version from Roblox and keeps the last ' + (pkg.Update.keep || 2) + '.' : 'Runs the app\'s own updater here.'}</p>
                    {pkg.LastLog && <p className="small muted"><code>{pkg.LastLog}</code></p>}
                    <div className="field">
                        <button className="primary" disabled={!!busy} onClick={() => run('Updating ' + pkg.Name, () => Go.UpdateNow(pkg.Dir))}>Update now</button>
                        <span className="spacer"/>
                        <span className="small muted">Automatic</span>
                        <select value={nearest(pkg.ScheduleMinutes)} disabled={!!busy}
                                onChange={e => run('Saving schedule', () => Go.SetSchedule(pkg.Dir, Number(e.target.value)), 'Schedule saved')}>
                            {SCHEDULES.map(([m, t]) => <option key={m} value={m}>{t}</option>)}
                        </select>
                    </div>
                </div>
            )}

            {hunt && <HuntCard pkg={pkg} st={hunt} busy={busy} run={run} close={closeHunt}/>}

            <div className="card">
                <div className="card-head"><LinkIcon/><span className="name">Folders</span><span className="spacer"/>
                    <button disabled={!!busy || (hunt && hunt.running)} onClick={startHunt}
                            title="Start the installed app and find the folders it uses"><SearchIcon/>Hunt…</button>
                    <button disabled={!!busy} onClick={() => run('Linking on this PC', () => Go.SetupHere(pkg.Dir))}>Link on this PC</button>
                    <button disabled={!!busy} onClick={() => run('Unlinking on this PC', () => Go.UnlinkHere(pkg.Dir))}>Unlink</button>
                </div>
                {(pkg.Folders || []).map((f, i) => (
                    <div key={i} className="field">
                        <span className={'pill' + (f.State === 1 ? ' ok' : '')}>{stateText[f.State]}</span>
                        <code className="hunt-path" title={f.Link + ' → ' + f.Data}>{pkg.Links[i]}</code>
                        <SessionBox value={(pkg.Session[i] || []).join(', ')} disabled={!!busy}
                                    save={v => run('Saving', () => Go.SetSession(pkg.Dir, i, v), 'Saved')}/>
                        <button className="icon-btn danger" title="Take out of the package (keeps its copy)" disabled={!!busy}
                                onClick={() => run('Removing folder', () => Go.RemoveFolder(pkg.Dir, i), 'Folder taken out')}><XIcon/></button>
                    </div>
                ))}
                <div className="field">
                    <input value={link} onChange={e => setLink(e.target.value)} placeholder="Where the app keeps a folder, such as %LOCALAPPDATA%\Game"/>
                    <button onClick={async () => { const d = await Go.PickFolder('A folder the app uses'); if (d) setLink(d) }}>Browse…</button>
                </div>
                <div className="field">
                    <input value={session} onChange={e => setSession(e.target.value)} placeholder="Subfolders to empty after each player (logins), comma-separated"/>
                    <button className="primary" disabled={!!busy || !link}
                            onClick={() => run('Adding folder', () => Go.AddFolder(pkg.Dir, link, session), 'Folder added').then(() => { setLink(''); setSession('') })}>Add folder</button>
                </div>
            </div>

            <div className="card">
                <div className="card-head"><KeyIcon/><span className="name">Registry</span></div>
                {(pkg.Registry || []).map(k => <code key={k} className="hunt-path">{k}</code>)}
                <div className="field">
                    <input value={key} onChange={e => setKey(e.target.value)} placeholder="HKCU\Software\Game"/>
                    <button className="primary" disabled={!!busy || !key}
                            onClick={() => run('Saving key', () => Go.AddRegistry(pkg.Dir, key), 'Key saved').then(() => setKey(''))}>Add key</button>
                </div>
            </div>

            <div className="card">
                <div className="card-head"><span className="name">Floor PCs (CCBoot)</span></div>
                <p className="small muted">Add this line to <code>C:\CCBoot\Cmd\public.txt</code> so each PC links the app when a customer logs on.
                    Use the drive letter the game disk has on the clients.</p>
                <div className="field">
                    <input readOnly value={pkg.CCBoot}/>
                    <button onClick={() => navigator.clipboard.writeText(pkg.CCBoot)}>Copy</button>
                </div>
            </div>
        </>
    )
}

// SessionBox edits a folder's subfolders emptied after each player; it
// saves when you leave it or press Enter.
function SessionBox({value, disabled, save}) {
    const [v, setV] = useState(value)
    useEffect(() => setV(value), [value])
    const commit = () => { if (v.trim() !== value) save(v) }
    return (
        <input className="session" value={v} disabled={disabled} placeholder="emptied after each player"
               title="Subfolders emptied after each player, such as LocalStorage (logins). Comma-separated."
               onChange={e => setV(e.target.value)} onBlur={commit}
               onKeyDown={e => { if (e.key === 'Enter') e.target.blur() }}/>
    )
}

// HuntCard shows what a hunt finds, and adds the ticked folders.
function HuntCard({pkg, st, busy, run, close}) {
    const [ticked, setTicked] = useState({})
    const [more, setMore] = useState(false)
    const folders = st.folders || []
    // Folders that are not noise and not added yet are ticked until unticked.
    const isTicked = f => !f.added && (ticked[f.real] ?? !f.noise)
    const main = folders.filter(f => !f.noise || f.added)
    const noise = folders.filter(f => f.noise && !f.added)
    const chosen = folders.filter(isTicked).map(f => f.real)
    const exeName = (st.exe || '').split('\\').pop()

    const row = f => (
        <label key={f.real} className={'hunt-row' + (f.added ? ' added' : '')} title={f.example}>
            <input type="checkbox" disabled={f.added} checked={f.added || isTicked(f)}
                   onChange={e => setTicked({...ticked, [f.real]: e.target.checked})}/>
            <span className="hunt-path">{f.path}</span>
            {f.added && <span className="pill ok">added</span>}
            {f.program && <span className="pill">program</span>}
            {f.changes > 0 && <span className="pill">{f.changes} changes</span>}
        </label>
    )

    return (
        <div className="card hunt">
            <div className="card-head"><SearchIcon/>
                <span className="name">Hunt: {exeName}</span><span className="spacer"/>
                {st.running
                    ? <><span className="tag">{st.live} running</span><button onClick={() => Go.StopHunt()}>Stop</button></>
                    : <button className="icon-btn" title="Close" onClick={close}><XIcon/></button>}
            </div>
            {st.running && <p className="small muted">Use the app the way players do: log in, change a setting, play a
                little. Then close it (also from the tray). Hunting stops a few seconds after its last program ends.
                Close other programs meanwhile: every change on the PC is counted.</p>}
            {st.error && <p className="error-text">{st.error}</p>}
            {st.programs && st.programs.length > 0 && <p className="small muted">Programs: {st.programs.join(', ')}</p>}
            <div className="hunt-list">
                {main.map(row)}
                {main.length === 0 && <p className="small muted">{st.running ? 'Nothing found yet…' : 'Nothing found.'}</p>}
                {more && noise.map(row)}
            </div>
            {noise.length > 0 && <button className="ghost hunt-more" onClick={() => setMore(!more)}>
                {more ? 'Hide' : 'Show'} {noise.length} more (most likely Windows or other programs)</button>}
            <div className="hunt-actions">
                <button className="primary" disabled={!!busy || chosen.length === 0}
                        onClick={() => run('Copying ' + chosen.length + ' folder(s) into the package',
                            () => Go.AddHuntFolders(pkg.Dir, chosen), 'Added ' + chosen.length + ' folder(s)')
                            .then(() => setTicked({}))}>
                    Add {chosen.length} folder{chosen.length === 1 ? '' : 's'} to the package
                </button>
            </div>
            <p className="small muted">The app's own program folder is already copied into the package when you made it.
                Add the folders it keeps elsewhere (settings, saves). Their files are copied in now.</p>
        </div>
    )
}

function nearest(m) {
    if (!m) return 0
    let best = SCHEDULES[1][0]
    for (const [v] of SCHEDULES) if (v && Math.abs(v - m) < Math.abs(best - m)) best = v
    return best
}
