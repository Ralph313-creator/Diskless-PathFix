import logo from './assets/logo.png'

// Brand is the app's name and maker, shown at the left of every tab's top bar.
export default function Brand() {
    return (
        <div className="brand">
            <img src={logo} alt=""/>
            <div className="brand-text">
                <h1>Diskless PathFix</h1>
                <span className="byline">Created by: Raptech</span>
            </div>
        </div>
    )
}
