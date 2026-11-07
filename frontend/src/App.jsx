import {useEffect, useState} from 'react'
import Brand from './Brand'
import * as Go from '../wailsjs/go/main/App'
import GameFixTab from './GameFix'
import HuntTab from './Hunt'
import {Guide, HelpIcon, Steps, store} from './Guide'
import PublicTxt, {copyText, FileIcon} from './PublicTxt'
import Update, {UpdateButton, useUpdate} from './Update'

const TABS = [
    ['portable', 'Portable apps', 'Apps that live on the game disk, such as Roblox'],
    ['gamefix', 'SymbolicLink', 'Folders games need on C:, copied to the game disk'],
    ['hunt', 'Hunt', 'Find the folders a program uses'],
]

// App is the window: a Portable apps tab, a SymbolicLink tab and the Hunt tool.
// All stay mounted, so a copy or hunt running in one keeps going while
// another is shown. The guide opens by itself the first time.
export default function App() {
    const [tab, setTab] = useState(() => store.get('tab') || 'portable')
    const [guide, setGuide] = useState(() => store.get('guide-seen') !== '1')
    const [lines, setLines] = useState(false) // the lines for public.txt are shown
    const [updates, setUpdates] = useState(false) // the Updates window is shown
    const update = useUpdate()
    const choose = t => {
        setTab(t)
        store.set('tab', t)
    }
    const closeGuide = () => {
        setGuide(false)
        store.set('guide-seen', '1')
    }
    const tabs = (
        <nav className="tabs">
            {TABS.map(([id, label, about]) => (
                <button key={id} className={'tab' + (tab === id ? ' active' : '')} title={about} onClick={() => choose(id)}>{label}</button>
            ))}
            <button className="tab help-tab" title="How Diskless PathFix works" onClick={() => setGuide(true)}><HelpIcon/>Help</button>
            <button className="tab help-tab" title="The lines to copy into CCBoot's C:\CCBoot\Cmd\public.txt" onClick={() => setLines(true)}><FileIcon/>public.txt</button>
            <UpdateButton info={update} onClick={() => setUpdates(true)}/>
        </nav>
    )
    const openLines = () => setLines(true)
    return <>
        <PortableTab tabs={tabs} hidden={tab !== 'portable'} openLines={openLines}/>
        <GameFixTab tabs={tabs} hidden={tab !== 'gamefix'} openLines={openLines}/>
        <HuntTab tabs={tabs} hidden={tab !== 'hunt'}/>
        {guide && <Guide onClose={closeGuide} openTab={t => { choose(t); closeGuide() }} openLines={() => { closeGuide(); openLines() }}/>}
        {lines && <PublicTxt onClose={() => setLines(false)}/>}
        {updates && <Update onClose={() => setUpdates(false)}/>}
    </>
}

// AppIcon shows a launcher's icon, or an empty square when it has none.
const AppIcon = ({src, size}) => src
    ? <img className="app-icon" src={src} width={size} height={size} alt=""/>
    : <span className="app-icon none" style={{width: size, height: size}}/>

const Icon = ({d}) => (
    <svg className="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"
         strokeLinecap="round" strokeLinejoin="round">{d}</svg>
)
const FolderIcon = () => <Icon d={<path d="M3 7a2 2 0 0 1 2-2h4l2 2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z"/>}/>
const KeyIcon = () => <Icon d={<><circle cx="8" cy="15" r="4"/><path d="m11 12 9-9M17 6l3 3"/></>}/>
const RefreshIcon = () => <Icon d={<><path d="M21 12a9 9 0 1 1-3-6.7L21 8"/><path d="M21 3v5h-5"/></>}/>
const XIcon = () => <Icon d={<path d="M18 6 6 18M6 6l12 12"/>}/>
const LinkIcon = () => <Icon d={<><path d="M10 13a5 5 0 0 0 7 0l3-3a5 5 0 0 0-7-7l-1 1"/><path d="M14 11a5 5 0 0 0-7 0l-3 3a5 5 0 0 0 7 7l1-1"/></>}/>

const SCHEDULES = [[0, 'Off'], [15, 'Every 15 minutes'], [30, 'Every 30 minutes'], [60, 'Every hour'], [360, 'Every 6 hours'], [1440 - 1, 'Once a day']]

const stateText = {0: 'not linked', 1: 'linked', 2: 'links elsewhere', 3: 'real folder (not linked)', 4: 'a file is in the way'}

function PortableTab({tabs, hidden, openLines}) {
    const [root, setRoot] = useState('')
    const [packages, setPackages] = useState([])
    const [selected, setSelected] = useState('')
    const [busy, setBusy] = useState('')
    const [toast, setToast] = useState(null)
    const [creating, setCreating] = useState('')

    const refresh = async () => {
        const list = (await Go.List()) || []
        setPackages(list)
        if (!list.find(p => p.Dir === selected) && list.length) setSelected(list[0].Dir)
    }
    useEffect(() => { Go.GetRoot().then(r => { setRoot(r); refresh() }) }, [])

    // run does one action with a busy label and shows how it went. It
    // resolves to true when the action worked.
    const run = async (label, fn, okText) => {
        setBusy(label)
        try {
            const out = await fn()
            setToast({ok: true, text: okText || (typeof out === 'string' && out ? lastLine(out) : 'Done')})
            return true
        } catch (e) {
            setToast({ok: false, text: String(e)})
            return false
        } finally {
            setBusy('')
            refresh()
        }
    }

    const pkg = packages.find(p => p.Dir === selected)
    const done = dir => { setCreating(''); setSelected(dir) }
    const chooseRoot = async () => { setRoot(await Go.ChooseRoot()); refresh() }
    const updatable = packages.filter(p => p.Update && !p.Error)

    const steps = [
        {
            title: 'Choose the apps folder', done: !!root,
            hint: <>Pick a folder on the game disk, such as <code>G:\Portable</code>. Each app gets its own folder in it.</>,
            action: <button className="primary" onClick={chooseRoot}>Choose apps folder…</button>,
        },
        {
            title: 'Add an app', done: packages.length > 0,
            hint: <><b>Add Roblox</b> downloads Roblox for you. <b>Add another app</b> copies an app that is installed
                on this PC. Both are at the bottom left.</>,
            action: <button className="primary" disabled={!!busy || !root} onClick={() => setCreating('roblox')}>Add Roblox</button>,
        },
        {
            title: 'Keep it updated', done: packages.length > 0 && updatable.every(p => p.ScheduleMinutes > 0),
            hint: <>Open the app, click <b>Update now</b>, then set <b>Automatic</b> so it updates by itself.
                Always update here on the server: floor PCs lose their changes when they restart.</>,
        },
        {
            title: 'Set up the floor PCs',
            hint: <>Click <b>public.txt</b> at the top: it has the lines for every app and SymbolicLink. Copy them,
                paste them into <code>C:\CCBoot\Cmd\public.txt</code>, save, and restart one floor PC to test.</>,
            action: <button className="primary" disabled={!packages.length} onClick={openLines}><FileIcon/>Lines for public.txt</button>,
        },
    ]

    return (
        <div className="app" style={hidden ? {display: 'none'} : undefined}>
            <div className="topbar">
                <Brand/>
                {tabs}
                <div className="main-folder" title="Folder on the game disk where the portable apps are kept">
                    <span>Apps folder</span>
                    <span className="path">{root || 'not chosen'}</span>
                    <button onClick={chooseRoot}>{root ? 'Change…' : 'Choose…'}</button>
                </div>
            </div>
            <div className="body">
                <div className="sidebar">
                    <div className="sidebar-title">Portable apps</div>
                    <div className="game-list">
                        {packages.map(p => (
                            <button key={p.Dir} className={'game-item' + (p.Dir === selected ? ' active' : '')}
                                    onClick={() => { setSelected(p.Dir); setCreating('') }}>
                                <span className="title-with-icon"><AppIcon src={p.Icon} size={18}/><span className="game-name">{p.Name}</span></span>
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
                    <Steps id="portable" steps={steps}/>
                    {creating === 'roblox' && <NewRoblox busy={busy} run={run} done={done} root={root} cancel={() => setCreating('')}/>}
                    {creating === 'app' && <NewApp busy={busy} run={run} done={done} root={root} cancel={() => setCreating('')}/>}
                    {!creating && !pkg && (!root ? (
                        <div className="empty">
                            <h2>Start here</h2>
                            <p className="muted">Choose a folder on the game disk to keep the portable apps in, such as <code>G:\Portable</code>.</p>
                            <button className="primary big" onClick={chooseRoot}>Choose apps folder…</button>
                        </div>
                    ) : (
                        <div className="empty">
                            <h2>No portable apps yet</h2>
                            <p className="muted">Add Roblox to download it straight from Roblox, or add an app that is installed on this PC.</p>
                            <div className="empty-actions">
                                <button className="primary big" disabled={!!busy} onClick={() => setCreating('roblox')}>Add Roblox</button>
                                <button className="big" disabled={!!busy} onClick={() => setCreating('app')}>Add another app</button>
                            </div>
                        </div>
                    ))}
                    {!creating && pkg && <Package pkg={pkg} busy={busy} run={run} openLines={openLines}/>}
                </div>
            </div>
            <div className="footer">
                <span>{busy ? busy + '…' : 'Update apps here on the server only: floor PCs throw away what they write to the game disk when they restart.'}</span>
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

function NewRoblox({busy, run, done, root, cancel}) {
    const [perUser, setPerUser] = useState(false)
    useEffect(() => { Go.RobloxPerUser().then(setPerUser) }, [])
    return (
        <div className="card">
            <div className="card-head"><h2>Add Roblox</h2></div>
            <p className="muted">Makes <code>{root}\Roblox</code> and downloads the newest Roblox Player into it straight from Roblox.
                When a customer closes Roblox, their login is wiped, so the next customer starts logged out.</p>
            <p className="small"><b>Where is Roblox installed on the floor PCs' image?</b> If you don't know, keep the first one.</p>
            <label className="watch"><input type="radio" checked={!perUser} onChange={() => setPerUser(false)}/>
                For everyone (C:\Program Files\Roblox), like this server</label>
            <label className="watch"><input type="radio" checked={perUser} onChange={() => setPerUser(true)}/>
                For one user (%LOCALAPPDATA%\Roblox)</label>
            <div className="hunt-actions">
                <button onClick={cancel} disabled={!!busy}>Cancel</button>
                <button className="primary" disabled={!!busy}
                        onClick={() => run('Downloading Roblox', () => Go.MakeRoblox(perUser)).then(ok => ok && done(root + '\\Roblox'))}>
                    Make and download
                </button>
            </div>
        </div>
    )
}

function NewApp({busy, run, done, root, cancel}) {
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
            <p className="muted">Copies an app that is installed on this PC into <code>{root}\{name || '<name>'}</code>.
                Its usual place becomes a link to that copy, so the app still thinks it is installed there, and
                everything it saves goes into the copy.</p>
            <p className="small muted">Click <b>Browse…</b> next to Program and pick the app's .exe. Folder and Name fill in by themselves.</p>
            <div className="field"><span className="field-label">Program</span>
                <input value={exe} onChange={e => setExe(e.target.value)} placeholder="C:\Games\Game\game.exe"/>
                <button onClick={pickExe}>Browse…</button></div>
            <div className="field"><span className="field-label">Folder</span>
                <input value={appDir} onChange={e => setAppDir(e.target.value)} placeholder="the app's whole folder"/>
                <button onClick={async () => { const d = await Go.PickFolder("The app's folder"); if (d) setAppDir(d) }}>Browse…</button></div>
            <div className="field"><span className="field-label">Name</span>
                <input value={name} onChange={e => setName(e.target.value)} placeholder="Game"/></div>
            <p className="small muted">Afterwards, add the other folders and registry keys the app keeps on C:, such as its
                settings in %LOCALAPPDATA%. The Hunt tab can find them for you.</p>
            <div className="hunt-actions">
                <button onClick={cancel} disabled={!!busy}>Cancel</button>
                <button className="primary" disabled={!!busy || !name}
                        onClick={() => run('Copying ' + name, () => Go.CreatePackage(name, appDir, exe), 'Made ' + name)
                            .then(ok => ok && done(root + '\\' + name))}>
                    Make package
                </button>
            </div>
        </div>
    )
}

function Package({pkg, busy, run, openLines}) {
    const [link, setLink] = useState('')
    const [session, setSession] = useState('')
    const [key, setKey] = useState('')
    const [removing, setRemoving] = useState(false)
    useEffect(() => setRemoving(false), [pkg.Dir])
    const removeButton = <button className="danger-btn" disabled={!!busy} onClick={() => setRemoving(true)}><XIcon/>Remove app</button>
    if (pkg.Error) {
        return <>
            <div className="card"><div className="card-head"><h2>{pkg.Name}</h2><span className="spacer"/>{removeButton}</div>
                <p className="error-text">{pkg.Error}</p></div>
            {removing && <RemoveCard pkg={pkg} busy={busy} run={run} cancel={() => setRemoving(false)}/>}
        </>
    }
    return (
        <>
            <div className="game-head">
                <div className="title-with-icon">
                    <AppIcon src={pkg.Icon} size={40}/>
                    <div>
                        <h2>{pkg.Name}</h2>
                        <p className="small muted">{pkg.Dir}</p>
                    </div>
                </div>
                <div className="actions">
                    <button onClick={() => Go.Open(pkg.Dir)}><FolderIcon/>Open</button>
                    {removeButton}
                </div>
            </div>
            <div className="field icon-row">
                <span className="small muted">{pkg.OwnIcon ? 'Launcher icon: chosen' : pkg.Icon ? "Launcher icon: the app's own" : 'Launcher icon: none yet (the app has none, or is not downloaded yet)'}</span>
                <span className="spacer"/>
                <button disabled={!!busy} title="Take the icon from a program, dll or .ico file"
                        onClick={() => run('Changing icon', async () => { if (await Go.ChangeIcon(pkg.Dir)) return 'Icon changed' })}>Change icon…</button>
                {pkg.OwnIcon && <button disabled={!!busy} onClick={() => run('Changing icon', () => Go.UseAppIcon(pkg.Dir), "Back to the app's icon")}>Use app's icon</button>}
            </div>

            {removing && <RemoveCard pkg={pkg} busy={busy} run={run} cancel={() => setRemoving(false)}/>}

            {pkg.Update && (
                <div className="card">
                    <div className="card-head"><RefreshIcon/><span className="name">Updates</span><span className="spacer"/>
                        {pkg.Version && <span className="pill ok">{pkg.Version}</span>}</div>
                    <p className="small muted">{pkg.Update.kind === 'roblox' ? 'Downloads each new Roblox version from Roblox and keeps the last ' + (pkg.Update.keep || 2) + '.' : 'Runs the app\'s own updater here.'}
                        {' '}Set <b>Automatic</b> so you never have to do it by hand.</p>
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

            <div className="card">
                <div className="card-head"><span className="name">Set up the floor PCs</span></div>
                <ol className="how">
                    <li>Copy this line:</li>
                </ol>
                <div className="field">
                    <input readOnly value={pkg.CCBoot}/>
                    <button className="primary" onClick={() => run('Copying', async () => {
                        if (!await copyText(pkg.CCBoot)) throw new Error('Could not copy. Select the line and press Ctrl+C.')
                    }, 'Line copied')}>Copy</button>
                </div>
                <ol className="how" start={2}>
                    <li>Open <code>C:\CCBoot\Cmd\public.txt</code> in Notepad, paste the line on its own line, and save.</li>
                    <li>Restart one floor PC to test: the app should start from its icon and from the website's Play button.</li>
                </ol>
                <p className="small muted">The line has this PC's drive letter. If the game disk has another letter on the floor PCs, change it in the line.</p>
                <div className="field">
                    <span className="small muted">More than one app, or SymbolicLink too? Copy all their lines at once:</span>
                    <button onClick={openLines}><FileIcon/>All lines for public.txt</button>
                </div>
            </div>

            <div className="section-title">What the app keeps</div>

            <div className="card">
                <div className="card-head"><LinkIcon/><span className="name">Folders</span><span className="spacer"/>
                    <button disabled={!!busy} title="Make the links on this PC, to test the app here" onClick={() => run('Linking on this PC', () => Go.SetupHere(pkg.Dir))}>Link on this PC</button>
                    <button disabled={!!busy} title="Remove the links on this PC and put the original folders back" onClick={() => run('Unlinking on this PC', () => Go.UnlinkHere(pkg.Dir))}>Unlink</button>
                </div>
                <p className="small muted">The app's folders are kept in the package. On a PC where the app is set up, each
                    one is linked from its usual place on C:. <b>"not linked" is normal on this server</b>: floor PCs make
                    the links when a customer logs on. Use <b>Link on this PC</b> only to test here.
                    Not sure which folders the app uses? Use the Hunt tab.</p>
                {(pkg.Folders || []).map((f, i) => (
                    <div key={i} className="field">
                        <span className={'pill' + (f.State === 1 ? ' ok' : '')}>{stateText[f.State]}</span>
                        <code className="hunt-path" title={f.Link + ' → ' + f.Data}>{pkg.Links[i]}</code>
                        <SessionBox value={(pkg.Session[i] || []).join(', ')} disabled={!!busy}
                                    save={v => run('Saving', () => Go.SetSession(pkg.Dir, i, v), 'Saved')}/>
                        <button className="icon-btn danger" title="Stop linking this folder (its files stay in the package). To remove the whole app, use Remove app at the top." disabled={!!busy}
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
                            onClick={() => run('Adding folder', () => Go.AddFolder(pkg.Dir, link, session), 'Folder added')
                                .then(ok => { if (ok) { setLink(''); setSession('') } })}>Add folder</button>
                </div>
            </div>

            <div className="card">
                <div className="card-head"><KeyIcon/><span className="name">Registry</span></div>
                <p className="small muted">Settings the app keeps in the Windows registry. They are saved into the package
                    when the app closes and loaded again next time, on any PC.</p>
                {(pkg.Registry || []).map(k => <code key={k} className="hunt-path">{k}</code>)}
                <div className="field">
                    <input value={key} onChange={e => setKey(e.target.value)} placeholder="HKCU\Software\Game"/>
                    <button className="primary" disabled={!!busy || !key}
                            onClick={() => run('Saving key', () => Go.AddRegistry(pkg.Dir, key), 'Key saved').then(ok => ok && setKey(''))}>Add key</button>
                </div>
            </div>
        </>
    )
}

// RemoveCard asks before removing an app, showing what will be closed.
function RemoveCard({pkg, busy, run, cancel}) {
    const [use, setUse] = useState(null)
    useEffect(() => { Go.InUse(pkg.Dir).then(setUse) }, [pkg.Dir])
    const programs = (use && use.Programs) || []
    const services = (use && use.Services) || []
    return (
        <div className="card danger">
            <div className="card-head"><span className="name">Remove {pkg.Name}?</span></div>
            <p>This deletes <code>{pkg.Dir}</code> with everything in it, removes its links and Play-button
                handler on this PC (the original folders come back), and turns off its automatic updates.</p>
            {use === null && <p className="small muted">Looking for programs that use it…</p>}
            {services.length > 0 && <p className="small">Windows services installed from it will be stopped and
                unregistered: <b>{services.join(', ')}</b>. The app installs them again when it is set up from its own installer.</p>}
            {programs.length > 0 && <p className="small">Running from it, will be closed:
                <b> {[...new Set(programs.map(p => p.Name))].join(', ')}</b></p>}
            <p className="small muted">Floor PCs that link it keep a dead link until they restart; remove its line from
                C:\CCBoot\Cmd\public.txt too.</p>
            <div className="hunt-actions">
                <button onClick={cancel} disabled={!!busy}>Cancel</button>
                <button className="danger-btn solid" disabled={!!busy || use === null}
                        onClick={() => run('Removing ' + pkg.Name, () => Go.RemoveApp(pkg.Dir), 'Removed ' + pkg.Name)}>
                    {programs.length + services.length > 0 ? 'Close them and remove' : 'Remove'}
                </button>
            </div>
        </div>
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

function nearest(m) {
    if (!m) return 0
    let best = SCHEDULES[1][0]
    for (const [v] of SCHEDULES) if (v && Math.abs(v - m) < Math.abs(best - m)) best = v
    return best
}
