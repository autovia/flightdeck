// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import React, {Component} from 'react';

function classNames(...classes) {
    return classes.filter(Boolean).join(' ')
  }

class Tabs extends Component {
    constructor(props) {
        super(props);
        this.state = {
            current: this.props.current
        }
    }

    onClick = (id) => {
        console.log("onClick: ", id);
        this.setState({
            current: id
        });
        this.props.onTabClick(id);
    };
    
    render() {
        return (
            <div>
            <div className="sm:hidden">
              <label htmlFor="tabs" className="sr-only">
                Select a tab
              </label>
            </div>
            <div className="hidden sm:block">
              <div className="border-b border-gray-200">
                <nav className="-mb-px flex space-x-8" aria-label="Tabs">
                  {this.props.tabs.map((tab) => (
                    <button
                      key={tab.name}
                      value={tab.id}
                      onClick={() => this.onClick(tab.id)}
                      className={classNames(
                        this.state.current == tab.id
                          ? 'border-indigo-500 text-indigo-600'
                          : 'border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700',
                        'whitespace-nowrap border-b-2 py-4 px-1 text-sm font-medium'
                      )}
                      aria-current={this.state.current == tab.id ? 'page' : undefined}
                    >
                      {tab.name}
                    </button>
                  ))}
                </nav>
              </div>
            </div>
          </div>
        );
    }
}

export default Tabs;