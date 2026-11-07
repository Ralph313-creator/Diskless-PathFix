import {useEffect, useState} from 'react'
import logo from './assets/logo.png'
import {GetVersion} from '../wailsjs/go/main/App'

// Brand is the app's name, maker and version, shown at the left of every tab's top bar.
export default function Brand() {
    const [version, setVersion] = useState('')
    useEffect(() => { GetVersion().then(setVersion) }, [])
    return (
        <div className="brand">
            <img src={logo} alt=""/>
            <div className="brand-text">
                <h1>Diskless PathFix</h1>
                <span className="byline">Created by: Raptech{version && ` · v${version}`}</span>
            </div>
        </div>
    )
}
