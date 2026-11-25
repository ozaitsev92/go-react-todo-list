import { Component } from "react";
import PropTypes from "prop-types";

// ErrorBoundary replaces a crashed render tree with a recoverable message.
export default class ErrorBoundary extends Component {
    state = { failed: false };

    static getDerivedStateFromError() {
        return { failed: true };
    }

    render() {
        if (this.state.failed) {
            return (
                <div role="alert" className="container mt-5 text-center">
                    <h1 className="h4">Something went wrong.</h1>
                    <button
                        type="button"
                        className="btn btn-primary"
                        onClick={() => window.location.reload()}
                    >
                        Reload
                    </button>
                </div>
            );
        }
        return this.props.children;
    }
}

ErrorBoundary.propTypes = {
    children: PropTypes.node.isRequired,
};
