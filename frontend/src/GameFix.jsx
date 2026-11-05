import {useEffect, useRef, useState} from 'react'
import Brand from './Brand'
import {
    AddFolder, AddGame, BrowseFolder, CCBootLine, ChooseFolder, ChooseMainFolder, CopyFolder, DeleteGame,
    GetMainFolder, GetServiceState, InstallService, ListGames, OpenFolder, RemoveFolder, SaveFolders,
    SetWatch, UninstallService
} from '../wailsjs/go/gamefix/App'
import {EventsOn} from '../wailsjs/runtime/runtime'
import logo from './assets/logo.png'
import {Steps} from './Guide'
import {copyText, FileIcon} from './PublicTxt'

const GAME_ITSELF = '.' // folder copied straight into the game folder; see GameItself in linkbat.go

function formatSize(bytes) {
    const units = ['B', 'KB', 'MB', 'GB', 'TB']
    let i = 0
    while (bytes >= 1024 && i < units.length - 1) {
        bytes /= 1024
        i++
    }
    return `${bytes.toFixed(i === 0 ? 0 : 1)} ${units[i]}`
}

const folderName = (game, folder) => folder.name === GAME_ITSELF ? game.name : folder.name
const copyPath = (game, folder) => folder.name === GAME_ITSELF ? game.folder : game.folder + '\\' + folder.name

const ICONS = {
    plus: 'M12 5v14M5 12h14',
    folder: 'M3 7a2 2 0 0 1 2-2h4l2 2h8a2 2 0 0 1 2 2v8a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2z',
    open: 'M14 4h6v6M20 4l-9 9M18 14v5a1 1 0 0 1-1 1H5a1 1 0 0 1-1-1V7a1 1 0 0 1 1-1h5',
    trash: 'M4 7h16M10 11v6M14 11v6M6 7l1 12a2 2 0 0 0 2 2h6a2 2 0 0 0 2-2l1-12M9 7V4h6v3',
    close: 'M6 6l12 12M18 6L6 18',
    search: 'M11 4a7 7 0 1 0 0 14a7 7 0 1 0 0-14M20 20l-4-4',
}

function Icon({name}) {
    return (
        <svg className="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"
             strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">
            <path d={ICONS[name]}/>
        </svg>
    )
}

function IconButton({icon, title, className = '', ...props}) {
    return (
        <button className={'icon-btn ' + className} title={title} aria-label={title} {...props}>
            <Icon name={icon}/>
        </button>
    )
}

// PathField is a path that saves when you leave it or press Enter.
function PathField({label, value, placeholder, hint, onChange, onCommit, onBrowse}) {
    return (
        <div className="field">
            <span className="field-label">{label}</span>
            <input
                value={value}
                placeholder={placeholder}
                title={hint}
                spellCheck={false}
                onChange={e => onChange(e.target.value)}
                onBlur={onCommit}
                onKeyDown={e => e.key === 'Enter' && e.currentTarget.blur()}
            />
            {onBrowse && <button className="ghost" onClick={onBrowse}>Browse</button>}
        </div>
    )
}

// GameFixTab is the SymbolicLink app, ported: copy game folders into a main folder
// for clients to link (GameFix.bat) or copy (GameFix Copy.bat).
export default function GameFixTab({tabs, hidden, openLines}) {
    const [mainFolder, setMainFolder] = useState('')
    const [loaded, setLoaded] = useState(false)
    const [games, setGames] = useState([])
    const [selected, setSelected] = useState('')
    const [newName, setNewName] = useState('')
    const [copying, setCopying] = useState(null) // {game, folder, file, done, total} while a copy runs
    const [toast, setToast] = useState(null) // {kind: 'ok' | 'error', text}
    const [edits, setEdits] = useState({}) // unsaved fields: {game: {folder: {source, link}}}
    const [service, setService] = useState('') // state of the Windows service; see serviceState in service.go
    const [serviceBusy, setServiceBusy] = useState(false)
    const [line, setLine] = useState('') // GameFix.bat's line for CCBoot's public.txt; see CCBootLine in ccboot.go
    const toastTimer = useRef()
    useEffect(() => { CCBootLine().then(setLine) }, [mainFolder])

    const game = games.find(g => g.name === selected) ?? games[0]

    const notify = (kind, text) => {
        clearTimeout(toastTimer.current)
        setToast({kind, text})
        if (kind === 'ok') toastTimer.current = setTimeout(() => setToast(null), 4000)
    }
    const fail = (err) => notify('error', String(err))

    const refresh = async () => {
        try {
            setGames(await ListGames())
        } catch (err) {
            fail(err)
        }
    }

    useEffect(() => {
        GetMainFolder().then(setMainFolder).then(refresh).finally(() => setLoaded(true))
        const offProgress = EventsOn('copy-progress', setCopying)
        // A watched folder was copied again because its source changed.
        const offWatch = EventsOn('watch-copied', ({game, folder, error}) => {
            setCopying(null)
            if (error) fail(`Copying ${folder === GAME_ITSELF ? game : folder}: ${error}`)
            refresh()
        })
        const checkService = () => GetServiceState().then(setService)
        checkService()
        const serviceTimer = setInterval(checkService, 5000)
        return () => {
            offProgress()
            offWatch()
            clearInterval(serviceTimer)
        }
    }, [])

    const toggleService = async () => {
        const installed = service !== 'not installed'
        if (installed && !confirm('Uninstall the SymbolicLink service?\n\nEnabled folders are then only copied while Diskless PathFix is open.')) return
        setServiceBusy(true)
        try {
            await (installed ? UninstallService() : InstallService())
            notify('ok', installed
                ? 'Service uninstalled. Enabled folders are copied while Diskless PathFix is open.'
                : 'Service installed. Enabled folders are now copied even when Diskless PathFix is closed.')
        } catch (err) {
            fail(err)
        }
        setService(await GetServiceState())
        setServiceBusy(false)
    }

    const valueOf = (g, f, field) => edits[g.name]?.[f.name]?.[field] ?? f[field]
    const setField = (g, f, field, value) => setEdits(e => ({
        ...e,
        [g.name]: {...e[g.name], [f.name]: {...e[g.name]?.[f.name], [field]: value}},
    }))
    const isDirty = (g) => g.folders.some(f => valueOf(g, f, 'source') !== f.source || valueOf(g, f, 'link') !== f.link)

    // save saves every edited field of a game, and reports whether it worked.
    const save = async (g, override = {}) => {
        const folders = Object.fromEntries(g.folders.map(f => [f.name, {
            source: override[f.name]?.source ?? valueOf(g, f, 'source'),
            link: override[f.name]?.link ?? valueOf(g, f, 'link'),
        }]))
        try {
            await SaveFolders(g.name, folders)
            setEdits(e => {
                const {[g.name]: _, ...rest} = e
                return rest
            })
            await refresh()
            return true
        } catch (err) {
            fail(err)
            return false
        }
    }
    const commit = (g) => isDirty(g) && save(g)

    const chooseMain = async () => {
        try {
            const path = await ChooseMainFolder()
            if (path) {
                setMainFolder(path)
                await refresh()
            }
        } catch (err) {
            fail(err)
        }
    }

    const addGame = async (e) => {
        e.preventDefault()
        try {
            const g = await AddGame(newName)
            setNewName('')
            setSelected(g.name)
            await refresh()
        } catch (err) {
            fail(err)
        }
    }

    const deleteGame = async (g) => {
        if (!confirm(`Delete "${g.name}"?\n\nThis permanently deletes ${g.folder} and everything in it.`)) return
        try {
            await DeleteGame(g.name)
            notify('ok', `Deleted ${g.name}`)
        } catch (err) {
            fail(err)
        }
        await refresh()
    }

    const addFolder = async (g) => {
        const from = await ChooseFolder(g.name)
        if (!from) return
        try {
            await AddFolder(g.name, from)
            notify('ok', 'Folder added. Click Copy to copy it into the main folder.')
        } catch (err) {
            fail(err)
        }
        await refresh()
    }

    const copyFolder = async (g, f) => {
        if (isDirty(g) && !await save(g)) return
        setCopying({game: g.name, folder: f.name, file: '', done: 0, total: 0})
        try {
            await CopyFolder(g.name, f.name)
            notify('ok', `Copied ${folderName(g, f)}`)
        } catch (err) {
            fail(err)
        }
        setCopying(null)
        await refresh()
    }

    const toggleWatch = async (g, f) => {
        if (isDirty(g) && !await save(g)) return
        try {
            await SetWatch(g.name, f.name, !f.watch)
            const name = folderName(g, f)
            if (f.watch) {
                notify('ok', `${name} is off. It is not in GameFix.bat or GameFix Copy.bat, and clients don't get it.`)
            } else if (!valueOf(g, f, 'source')) {
                notify('ok', `${name} is on. Clients get it once it is copied. Set its source to copy it.`)
            } else {
                notify('ok', `${name} is on. Clients get it, and it is copied again whenever its source changes` +
                    (service === 'running' ? ', even when Diskless PathFix is closed.' : ', while Diskless PathFix is open. Install the service to keep copying it when SymbolicLink is closed.'))
            }
        } catch (err) {
            fail(err)
        }
        await refresh()
    }

    const removeFolder = async (g, f) => {
        const question = f.name === GAME_ITSELF
            ? `Remove "${g.name}"?\n\nThis deletes the copied game files in ${g.folder}` +
            (g.folders.length > 1 ? ', but keeps the copies of its other folders.' : '.') + ' The source is not touched.'
            : `Remove "${f.name}"?\n\nThis deletes its copy in ${copyPath(g, f)}. The source is not touched.`
        if (!confirm(question)) return
        try {
            await RemoveFolder(g.name, f.name)
        } catch (err) {
            fail(err)
        }
        await refresh()
    }

    const browse = async (g, f, field, title) => {
        const path = await BrowseFolder(title, valueOf(g, f, field))
        if (path) await save(g, {[f.name]: {[field]: path}})
    }

    if (!loaded) return null

    if (!mainFolder) {
        return (
            <div className="welcome" style={hidden ? {display: 'none'} : undefined}>
                {tabs}
                <img className="welcome-logo" src={logo} alt=""/>
                <h1>SymbolicLink</h1>
                <p>Some games need a folder on C:, such as Riot Games in C:\ProgramData. SymbolicLink copies those folders
                    to the game disk, and the floor PCs link to the copies.</p>
                <p>Start by choosing a main folder on the game disk to keep the copies in, such as <code>G:\GameFix</code>.</p>
                <button className="primary big" onClick={chooseMain}>Choose main folder</button>
                {toast && <Toast toast={toast} onClose={() => setToast(null)}/>}
            </div>
        )
    }

    const busy = !!copying
    const percent = copying?.total > 0 ? Math.min(100, copying.done / copying.total * 100) : 0
    const copyLine = async () => {
        if (await copyText(line)) notify('ok', 'Line copied. Paste it into C:\\CCBoot\\Cmd\\public.txt and save.')
        else fail('Could not copy. Select the line and press Ctrl+C.')
    }
    const steps = [
        {
            title: 'Add a game', done: games.length > 0,
            hint: <>Type the game's name in the box at the bottom left, such as <b>Riot</b>, and press Enter.</>,
        },
        {
            title: 'Add its folders', done: games.some(g => g.folders.length > 0),
            hint: <>Open the game and click <b>Add folder</b>. Choose a folder the game needs on C:, such
                as <code>C:\ProgramData\Riot Games</code>. Not sure which? The <b>Hunt</b> tab finds them.</>,
        },
        {
            title: 'Copy and enable', done: games.some(g => g.folders.some(f => f.copied && f.watch)),
            hint: <>Click <b>Copy</b>. Check <b>Link</b>: it is where the floor PCs need the folder, usually the same
                place as Source. Then tick <b>Enable</b>.</>,
        },
        {
            title: 'Run GameFix.bat on the floor PCs',
            hint: <>On a floor PC, <code>GameFix.bat</code> links every enabled folder to its copy. It needs admin
                rights, so paste this line into <code>C:\CCBoot\Cmd\public.txt</code> (in place of any older GameFix.bat line):
                <code className="line">{line}</code> With <b>-service</b>, CCBoot runs it as SYSTEM, with admin rights and
                nothing to click, each time a floor PC starts. Save, restart one floor PC, and test the game.</>,
            action: <div className="step-buttons">
                <button className="primary" disabled={!line} onClick={copyLine}>Copy line</button>
                <button onClick={openLines}><FileIcon/>All lines</button>
            </div>,
        },
    ]

    return (
        <div className="app" style={hidden ? {display: 'none'} : undefined}>
            <header className="topbar">
                <Brand/>
                {tabs}
                <div className="main-folder" title="Main folder">
                    <Icon name="folder"/>
                    <span className="path">{mainFolder}</span>
                    <button className="ghost" onClick={chooseMain} disabled={busy}>Change</button>
                    <IconButton icon="open" title="Open main folder" onClick={() => OpenFolder(mainFolder)}/>
                </div>
            </header>

            <div className="body">
                <aside className="sidebar">
                    <div className="sidebar-title">Games</div>
                    <nav className="game-list">
                        {games.map(g => (
                            <button
                                key={g.name}
                                className={'game-item' + (g === game ? ' active' : '')}
                                onClick={() => setSelected(g.name)}
                            >
                                <span className="game-name">{g.name}</span>
                                <span className="count">{g.folders.filter(f => f.copied).length}/{g.folders.length}</span>
                            </button>
                        ))}
                    </nav>
                    <form className="new-game" onSubmit={addGame}>
                        <input value={newName} onChange={e => setNewName(e.target.value)} placeholder="New game name"/>
                        <IconButton icon="plus" title="Add game" type="submit" className="primary" disabled={!newName.trim()}/>
                    </form>
                </aside>

                <main className="content">
                    <Steps id="gamefix" steps={steps}/>
                    {!game ? (
                        <div className="empty">
                            <h2>No games yet</h2>
                            <p className="muted">Type a game name in the box at the bottom left, such as Riot, and press Enter.</p>
                        </div>
                    ) : (
                        <>
                            <div className="game-head">
                                <div>
                                    <h2>{game.name}</h2>
                                    <div className="muted small">{game.folder}</div>
                                </div>
                                <div className="actions">
                                    <button className="primary" onClick={() => addFolder(game)} disabled={busy}>
                                        <Icon name="plus"/> Add folder
                                    </button>
                                    <IconButton icon="open" title="Open game folder" onClick={() => OpenFolder(game.folder)}/>
                                    <IconButton icon="trash" title="Delete game" className="danger"
                                                onClick={() => deleteGame(game)} disabled={busy}/>
                                </div>
                            </div>

                            {game.folders.length === 0 ? (
                                <div className="empty">
                                    <h2>No folders yet</h2>
                                    <p className="muted">Click <b>Add folder</b> and choose a folder the game needs on C:,
                                        such as its data, save or config folder. Not sure which? The Hunt tab finds them while the game runs.</p>
                                </div>
                            ) : game.folders.map(f => (
                                <section className="card" key={f.name}>
                                    <div className="card-head">
                                        <Icon name="folder"/>
                                        <span className="name">{folderName(game, f)}</span>
                                        {f.name === GAME_ITSELF && <span className="tag">game folder</span>}
                                        <span className={'pill' + (f.copied ? ' ok' : '')}>{f.copied ? 'Copied' : 'Not copied'}</span>
                                        <span className="spacer"/>
                                        <label className={'watch' + (f.watch ? ' on' : '')}
                                               title={'On: clients get this folder (GameFix.bat and GameFix Copy.bat),\n' +
                                                   'and it is copied again whenever its source changes: while SymbolicLink is open,\n' +
                                                   'or all the time if the service is installed.\n' +
                                                   'Off: clients don\'t get it.'}>
                                            <input type="checkbox" checked={f.watch} onChange={() => toggleWatch(game, f)}/>
                                            Enable
                                        </label>
                                        <button className={f.copied ? '' : 'primary'} onClick={() => copyFolder(game, f)} disabled={busy}>
                                            {f.copied ? 'Copy again' : 'Copy'}
                                        </button>
                                        <IconButton icon="open" title="Open copy" onClick={() => OpenFolder(copyPath(game, f))}
                                                    disabled={!f.copied}/>
                                        <IconButton icon="trash" title="Remove folder" className="danger"
                                                    onClick={() => removeFolder(game, f)} disabled={busy}/>
                                    </div>

                                    {copying?.game === game.name && copying?.folder === f.name ? (
                                        <div className="progress">
                                            <div className="bar"><div style={{width: `${percent}%`}}/></div>
                                            <div className="progress-text">
                                                <span className="file">{copying.file || 'Starting…'}</span>
                                                <span>{formatSize(copying.done)}{copying.total > 0 && ` / ${formatSize(copying.total)}`}</span>
                                            </div>
                                        </div>
                                    ) : (
                                        <>
                                            <PathField
                                                label="Source"
                                                value={valueOf(game, f, 'source')}
                                                placeholder="Where to copy it from"
                                                hint="User folders become variables such as %USERPROFILE% or %APPDATA%."
                                                onChange={v => setField(game, f, 'source', v)}
                                                onCommit={() => commit(game)}
                                                onBrowse={() => browse(game, f, 'source', 'Choose the source folder')}
                                            />
                                            <PathField
                                                label="Link"
                                                value={valueOf(game, f, 'link')}
                                                placeholder="Where clients get it, e.g. C:\Games\Game or %USERPROFILE%\Documents\Game"
                                                hint={'Where clients get this folder: GameFix.bat links it there,\n' +
                                                    'and GameFix Copy.bat copies it there.\n' +
                                                    'User folders become variables such as %USERPROFILE% or %APPDATA%,\n' +
                                                    'so it works for whoever is logged on to the client.'}
                                                onChange={v => setField(game, f, 'link', v)}
                                                onCommit={() => commit(game)}
                                                onBrowse={() => browse(game, f, 'link', 'Choose the link path')}
                                            />
                                            {f.userFolder && (
                                                <p className="small warn-text">This link path is in the customer's own folder
                                                    ({f.link.split('\\')[0]}). CCBoot runs GameFix.bat as SYSTEM, before anyone
                                                    logs on, so it skips this folder there. Run GameFix.bat by hand on a floor
                                                    PC to link it.</p>
                                            )}
                                        </>
                                    )}
                                </section>
                            ))}
                        </>
                    )}
                </main>
            </div>

            <footer className="footer">
                <span>
                    Floor PCs get the enabled folders by running <b>GameFix.bat</b> from the main folder.
                    <b> GameFix Copy.bat</b> gives each PC its own copy instead of a link.
                </span>
                <span className="service" title={'The SymbolicLink service copies the enabled folders\n' +
                    'whenever their source changes, even when this app is closed.'}>
                    <span className={'dot ' + (service === 'running' ? 'ok' : '')}/>
                    Service: {service === 'running' ? 'watching' : service || '…'}
                    <button className="ghost" onClick={toggleService} disabled={serviceBusy || !service}>
                        {service === 'not installed' ? 'Install' : 'Uninstall'}
                    </button>
                </span>
            </footer>

            {toast && <Toast toast={toast} onClose={() => setToast(null)}/>}
        </div>
    )
}


function Toast({toast, onClose}) {
    return (
        <div className={'toast ' + toast.kind} role="status">
            <span>{toast.text}</span>
            <IconButton icon="close" title="Close" onClick={onClose}/>
        </div>
    )
}
