import {useEffect, useState} from 'react'
import logo from './assets/logo.png'

// store is localStorage that never throws (it can be off or empty).
export const store = {
    get: k => { try { return localStorage.getItem(k) } catch { return null } },
    set: (k, v) => { try { localStorage.setItem(k, v) } catch {} },
}

const Icon = ({d}) => (
    <svg className="icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"
         strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">{d}</svg>
)
export const HelpIcon = () => <Icon d={<><circle cx="12" cy="12" r="9"/><path d="M9.5 9a2.5 2.5 0 1 1 3.5 2.3c-.6.3-1 .9-1 1.6V14M12 17.5v.01"/></>}/>
const XIcon = () => <Icon d={<path d="M18 6 6 18M6 6l12 12"/>}/>

// Steps is the "Getting started" box at the top of a tab. Each step is
// {title, hint, done, action}: done true or false ticks itself from what
// the app sees; done undefined is a step the app can't check (such as
// editing CCBoot's public.txt). The next step to do is highlighted and its
// help shown; clicking a step shows its help. "Hide steps" is remembered
// per tab (id).
export function Steps({id, steps}) {
    const [hidden, setHidden] = useState(() => store.get('steps-hidden:' + id) === '1')
    const [picked, setPicked] = useState(null)
    const next = steps.findIndex(s => s.done === false)
    useEffect(() => setPicked(null), [next]) // moving on shows the new next step
    const hide = h => {
        setHidden(h)
        store.set('steps-hidden:' + id, h ? '1' : '0')
    }

    if (hidden) {
        return (
            <div className="steps slim">
                <span className="steps-title">Getting started</span>
                <span className="small muted">{next >= 0 ? `Next: ${next + 1}. ${steps[next].title}` : 'All steps done'}</span>
                <span className="spacer"/>
                <button className="ghost" onClick={() => hide(false)}>Show steps</button>
            </div>
        )
    }

    const shown = picked ?? (next >= 0 ? next : steps.length - 1)
    const s = steps[shown]
    return (
        <div className="steps">
            <div className="steps-head">
                <span className="steps-title">Getting started</span>
                <span className="spacer"/>
                <button className="ghost" onClick={() => hide(true)}>Hide steps</button>
            </div>
            <ol className="step-list">
                {steps.map((st, i) => (
                    <li key={i}>
                        <button className={'step' + (i === shown ? ' shown' : '') + (st.done ? ' done' : '') +
                            (i === next ? ' next' : '') + (st.done === undefined ? ' info' : '')}
                                onClick={() => setPicked(i)}>
                            <span className="step-num">{st.done ? '✓' : i + 1}</span>
                            {st.title}
                        </button>
                    </li>
                ))}
            </ol>
            <div className="step-help">
                <div className="step-text">
                    <div className="step-title">
                        {shown + 1}. {s.title}
                        {s.done && <span className="pill ok">done</span>}
                        {s.done === undefined && <span className="pill">you do this</span>}
                    </div>
                    <p>{s.hint}</p>
                </div>
                {s.action && !s.done && <div className="step-action">{s.action}</div>}
            </div>
        </div>
    )
}

// Guide explains the whole program in plain words. It opens by itself the
// first time, and from Help next to the tabs. openTab(id) goes to a tab.
export function Guide({onClose, openTab, openLines}) {
    useEffect(() => {
        const key = e => e.key === 'Escape' && onClose()
        window.addEventListener('keydown', key)
        return () => window.removeEventListener('keydown', key)
    }, [onClose])
    return (
        <div className="overlay" onClick={onClose}>
            <div className="guide" role="dialog" aria-label="How RapTech Portable works" onClick={e => e.stopPropagation()}>
                <div className="guide-head">
                    <img src={logo} alt=""/>
                    <h2>How RapTech Portable works</h2>
                    <span className="spacer"/>
                    <button className="icon-btn" title="Close" aria-label="Close" onClick={onClose}><XIcon/></button>
                </div>
                <p className="guide-lead">
                    Floor PCs forget everything when they restart, and their C: drive is only as new as the last image
                    update. RapTech Portable keeps apps and game files on the <b>game disk (G:)</b> instead, so every
                    floor PC gets the newest version and the files it needs, without updating the image.
                </p>

                <h3>Which tab do I use?</h3>
                <div className="guide-cards">
                    <button className="guide-card" onClick={() => openTab('portable')}>
                        <b>Portable apps</b>
                        <span>For apps that update a lot, like <b>Roblox</b>. The whole app lives in one folder on G:.
                            Update it here and every floor PC has the new version.</span>
                        <span className="guide-go">Open Portable apps →</span>
                    </button>
                    <button className="guide-card" onClick={() => openTab('gamefix')}>
                        <b>GameFix</b>
                        <span>For games that need a folder on C:, like <b>Riot Games</b> in C:\ProgramData. GameFix copies
                            that folder to G:, and floor PCs link to it with GameFix.bat.</span>
                        <span className="guide-go">Open GameFix →</span>
                    </button>
                    <button className="guide-card" onClick={() => openTab('hunt')}>
                        <b>Hunt</b>
                        <span>Don't know which folders a program uses? Start it from Hunt, use it for a minute, close it.
                            Hunt lists the folders so you can add them.</span>
                        <span className="guide-go">Open Hunt →</span>
                    </button>
                </div>

                <h3>Words you will see</h3>
                <dl className="guide-words">
                    <dt>Link</dt>
                    <dd>A folder on C: that really opens a folder on G:. The app thinks its files are in the usual
                        place, but they are on the game disk.</dd>
                    <dt>Package</dt>
                    <dd>One app's folder in the apps folder, such as <code>G:\Portable\Roblox</code>. It holds the app,
                        its settings, and its launcher (<code>Roblox.exe</code>) that sets everything up and starts it.</dd>
                    <dt>Main folder</dt>
                    <dd>Where GameFix keeps the copied game folders, such as <code>G:\GameFix</code>.</dd>
                    <dt>Server</dt>
                    <dd>This PC, the CCBoot server. Floor PCs are the diskless customer PCs.</dd>
                </dl>

                <h3>Two rules</h3>
                <ol className="guide-rules">
                    <li><b>Update and copy here on the server.</b> Floor PCs throw away what they write to G: when they
                        restart, so only changes made here stay.</li>
                    <li><b>Floor PCs only do what CCBoot tells them.</b> Click <b>public.txt</b> next to the tabs, copy
                        the lines, and paste them into <code>C:\CCBoot\Cmd\public.txt</code>. Then restart one floor PC to
                        test. <button className="link-btn" onClick={openLines}>Show the lines</button></li>
                </ol>

                <div className="guide-foot">
                    <span className="small muted">Open this again any time with <b>Help</b> next to the tabs.</span>
                    <span className="spacer"/>
                    <button className="primary big" onClick={onClose}>Got it</button>
                </div>
            </div>
        </div>
    )
}
