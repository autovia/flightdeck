// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import React, {Component} from 'react';
import { XMarkIcon } from '@heroicons/react/24/outline'

class ListView extends Component {
  constructor(props) {
    super(props);
    this.state = {
        open: true,
        data: []
    }
  }

  componentDidMount() {
    fetch('/api/v1/' + this.props.meta.kind)
    .then(res => res.json())
    .then(d => {
      this.setState((state, props) => ({
        data: d,
      }));
    });
  }

  render() {
    return (
          <div className="flex h-full flex-col overflow-y-scroll bg-white bg-opacity-80 py-6 shadow-xl">
            <div className="px-4 sm:px-6">
              <div className="flex items-start justify-between">
                <div className="">
                  <div className="font-bold text-xl leading-6 text-gray-900 border-b-4 p-4">
                    {this.props.meta.label}
                  </div>
                </div>
                <div className="ml-3 flex h-7 items-center">
                  <button
                    type="button"
                    className="rounded-md bg-white text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                    onClick={this.props.close}
                  >
                    <span className="sr-only">Close panel</span>
                    <XMarkIcon className="h-6 w-6" aria-hidden="true" />
                  </button>
                </div>
              </div>
            </div>
            <div className="relative mt-6 flex-1 px-4 sm:px-6">
              <table className="table-auto w-300">
                <thead className="bg-white">
                  <tr>
                    <th scope="col" className="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900 sm:pl-3">
                      Name
                    </th>
                    <th scope="col" className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">
                      Created at
                    </th>
                  </tr>
                </thead>
                  <tbody className="bg-white">
                  {this.state.data.map((i) => (
                    <tr key={i.metadata.uid} className="border-t border-gray-200">
                      <td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-3">
                      {i.metadata.name}
                      </td>
                      <td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-3">
                      {i.metadata.creationTimestamp}
                      </td>     
                    </tr>
                  ))}
                  </tbody>
              </table>
            </div> 
          </div>
    );
  }
}

export default ListView;