import {useEffect, useState} from 'react'
import * as Go from '../wailsjs/go/main/App'
import * as GF from '../wailsjs/go/gamefix/App'
import {EventsOn} from '../wailsjs/runtime/runtime'
import {Steps} from './Guide'

// HuntTab is hunt mode, one tool for the whole program: hunt any program,
// then add the folders it found to a portable app or to a GameFix game.

const Icon = ({d}) => (
    <svg className="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"
         strokeLinecap="round" strokeLinejoin="round">{d}</svg>
)
const SearchIcon = () => <Icon d={<><circle cx="11" cy="11" r="7"/><path d="m21 21-4.3-4.3"/></>}/>
const XIcon = () => <Icon d={<path d="M18 6 6 18M6 6l12 12"/>}/>

const norm = p => (p || '').toLowerCase().replace(/[\\/]+$/, '')

export default function HuntTab({tabs, hidden, logo}) {
    const [st, setSt] = useState(null) // the hunt's state; see hunt.State
    const [targets, setTargets] = useState([]) // [{id, label}]
    const [target, setTarget] = useState('')   // "p:<package folder>" or "g:<game name>"
    const [have, setHave] = useState([])       // folders the target already has
    const [ticked, setTicked] = useState({})   // folder -> ticked, where the user changed it
    const [more, setMore] = useState(false)
    const [busy, setBusy] = useState(false)
    const [toast, setToast] = useState(null)
    const [addedOnce, setAddedOnce] = useState(false) // folders were added from this hunt

    useEffect(() => EventsOn('hunt', setSt), [])

    // The apps and games can change in the other tabs: list them again
    // whenever this tab is shown.
    const loadTargets = async () => {
        const list = []
        for (const p of (await Go.List()) || []) {
            if (!p.Error) list.push({id: 'p:' + p.Dir, label: 'Portable app: ' + p.Name})
        }
        try {
            for (const g of (await GF.ListGames()) || []) list.push({id: 'g:' + g.name, label: 'GameFix game: ' + g.name})
        } catch { /* no GameFix main folder yet */ }
        setTargets(list)
        if (!list.find(t => t.id === target)) setTarget(list.length ? list[0].id : '')
    }
    useEffect(() => { if (!hidden) loadTargets() }, [hidden])

    const loadHave = async (t = target) => {
        if (!t) return setHave([])
        try {
            const list = t.startsWith('p:') ? await Go.HuntAdded(t.slice(2)) : await GF.GameSources(t.slice(2))
            setHave((list || []).map(norm))
        } catch { setHave([]) }
    }
    useEffect(() => { loadHave() }, [target])

    const start = async () => {
        try {
            const exe = await Go.StartHunt()
            if (exe) {
                setTicked({})
                setAddedOnce(false)
                setSt({exe, running: true, programs: [], folders: [], live: 0})
            }
        } catch (e) { setToast({ok: false, text: String(e)}) }
    }

    const folders = (st && st.folders) || []
    const added = f => have.includes(norm(f.real))
    const isTicked = f => !added(f) && (ticked[f.real] ?? !f.noise)
    const main = folders.filter(f => !f.noise || added(f))
    const noise = folders.filter(f => f.noise && !added(f))
    const chosen = folders.filter(isTicked).map(f => f.real)
    const exeName = st ? (st.exe || '').split('\\').pop() : ''

    const add = async () => {
        setBusy(true)
        try {
            const name = t => t.slice(2)
            const n = target.startsWith('p:')
                ? await Go.AddHuntFolders(name(target), chosen)
                : await GF.AddHuntFolders(name(target), chosen)
            const where = (targets.find(t => t.id === target) || {}).label || ''
            setToast({ok: true, text: `Added ${n} folder${n === 1 ? '' : 's'} to ${where.replace(/^[^:]*: /, '')}.` +
                    (target.startsWith('g:') ? ' Check their link paths in the GameFix tab, then Copy and Enable them.' : ' Their files were copied into the package.')})
            setTicked({})
            setAddedOnce(true)
        } catch (e) { setToast({ok: false, text: String(e)}) }
        setBusy(false)
        loadHave()
    }

    const steps = [
        {
            title: 'Choose the program', done: !!st,
            hint: <>Close other programs first: every change on the PC is counted. Then click <b>Choose program…</b> below
                and pick the program that starts the app or game. It starts by itself.</>,
        },
        {
            title: 'Use it, then close it', done: !!st && !st.running,
            hint: <>Use it the way customers do: log in, change a setting, play a little. Then close it (also from the
                tray). Hunting stops a few seconds after it closes.</>,
        },
        {
            title: 'Add the folders', done: addedOnce,
            hint: <>The folders it used are ticked. Under <b>Add to</b>, choose the portable app or GameFix game they
                belong to, then click <b>Add</b>.</>,
        },
    ]

    const row = f => (
        <label key={f.real} className={'hunt-row' + (added(f) ? ' added' : '')} title={f.example}>
            <input type="checkbox" disabled={added(f)} checked={added(f) || isTicked(f)}
                   onChange={e => setTicked({...ticked, [f.real]: e.target.checked})}/>
            <span className="hunt-path">{f.path}</span>
            {added(f) && <span className="pill ok">already in it</span>}
            {f.program && <span className="pill">program</span>}
            {f.changes > 0 && <span className="pill">{f.changes.toLocaleString()} change{f.changes === 1 ? '' : 's'}</span>}
        </label>
    )

    return (
        <div className="app" style={hidden ? {display: 'none'} : undefined}>
            <div className="topbar">
                <div className="brand"><img src={logo} alt=""/><h1>RapTech Portable</h1></div>
                {tabs}
            </div>
            <div className="body">
                <div className="content">
                    <Steps id="hunt" steps={steps}/>
                    <div className="card">
                        <div className="card-head"><SearchIcon/><span className="name">Hunt: find the folders a program uses</span>
                            <span className="spacer"/>
                            {st && st.running
                                ? <><span className="pill ok">{st.live > 0 ? `${st.live} running` : 'Starting…'}</span>
                                    <button onClick={() => Go.StopHunt()} title="Stop hunting. The program keeps running.">Stop</button></>
                                : <button className="primary" onClick={start}><SearchIcon/>Choose program…</button>}
                        </div>
                        <p className="small muted">Hunt starts a program and watches it. It finds every folder the program uses and
                            every folder where files change, so you can add them to a portable app or a GameFix game.
                            Use the installed program, not a copy.</p>
                    </div>

                    {st && (
                        <div className="card hunt">
                            <div className="card-head">
                                <span className="name">{st.running ? 'Hunting' : 'Found'}: {exeName}</span>
                                <span className="spacer"/>
                                {!st.running && <button className="icon-btn" title="Close" onClick={() => setSt(null)}><XIcon/></button>}
                            </div>
                            {st.error && <p className="error-text">{st.error}</p>}
                            {st.programs && st.programs.length > 0 && <p className="small muted">Programs: {st.programs.join(', ')}</p>}
                            <div className="hunt-list">
                                {main.map(row)}
                                {main.length === 0 && <p className="small muted">{st.running ? 'Nothing found yet…' : 'Nothing found.'}</p>}
                                {more && noise.map(row)}
                            </div>
                            {noise.length > 0 && <button className="ghost hunt-more" onClick={() => setMore(!more)}>
                                {more ? 'Hide' : 'Show'} {noise.length} more (most likely Windows or other programs)</button>}
                            <div className="field">
                                <span className="small muted">Add to</span>
                                <select value={target} onChange={e => setTarget(e.target.value)} disabled={busy}>
                                    {targets.length === 0 && <option value="">Make a portable app or a GameFix game first</option>}
                                    {targets.map(t => <option key={t.id} value={t.id}>{t.label}</option>)}
                                </select>
                                <span className="spacer"/>
                                <button className="primary" disabled={busy || !target || chosen.length === 0} onClick={add}>
                                    Add {chosen.length} folder{chosen.length === 1 ? '' : 's'}
                                </button>
                            </div>
                        </div>
                    )}
                </div>
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
