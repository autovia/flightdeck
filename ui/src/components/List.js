// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import React, {Component} from 'react';
import { ArrowLeftIcon } from '@heroicons/react/24/outline'

class List extends Component {
  constructor(props) {
    super(props);
    this.state = {
        open: true,
        data: []
    }
  }

  componentDidMount() {
    this.getList();
  }

  getList() {
    fetch('/api/v1/list/' + this.props.meta.kind)
    .then(res => res.json())
    .then(d => {
      if (d && d.length > 0) {
        this.setState((state, props) => ({
          data: d,
        }));
      } else {
        this.setState((state, props) => ({
          data: [],
        }));
      }
    });
  }

  componentDidUpdate(prevProps, prevState) {
    if(prevState !== this.state) {
      console.log('List componentDidUpdate state: ', this.state);
    }
    if(prevProps !== this.props) {
      console.log('List componentDidUpdate props: ', this.props);
      this.getList();
    }
  }

  render() {
    return (
          <div className="flex h-full flex-col overflow-y-scroll bg-white bg-opacity-80 py-6 shadow-xl">
            <div className="px-4 sm:px-6">
              <div className="flex items-start justify-between">
                <div className="">
                  <div className="font-bold text-xl leading-6 text-gray-900 border-b-4 p-4">
                    <button
                      type="button"
                      className="pr-4 rounded-md bg-white text-gray-400 hover:text-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                      onClick={this.props.close}
                    >
                      <span className="sr-only">Back</span>
                      <ArrowLeftIcon className="h-6 w-6" aria-hidden="true" />
                    </button>
                    {this.props.meta.label}
                  </div>
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
                    <tr key={i.uid} className="border-t border-gray-200">
                      <td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-3">
                      {i.name}
                      </td>
                      <td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900 sm:pl-3">
                      {i.creationTimestamp}
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

export default List;