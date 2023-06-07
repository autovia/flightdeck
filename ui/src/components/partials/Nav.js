// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import { Fragment, Component } from 'react'
import { Disclosure, Menu, Transition } from '@headlessui/react'
import { UserCircleIcon, MagnifyingGlassIcon } from '@heroicons/react/20/solid'
import MenuSelect from './MenuSelect';
import MenuSelectResource from './MenuSelectResource';
import MenuSelectCluster from './MenuSelectCluster';
import * as k8s from '../utils/K8s';

const userNavigation = [
  { name: 'Sign out', href: '#' },
]

class Nav extends Component {

  constructor(props) {
    super(props);
    this.state = {
      namespace: props.namespace,
      namespaces: [],
      resource: {}, 
      clusterResource: {}
    }
  }

  componentDidMount() {
    console.log("Nav props", this.props);
    fetch('/api/v1/namespaces')
    .then(res => res.json())
    .then(d => {
      console.log('/api/v1/namespaces', d);
      const namespaces = d.nodes.map((n) => ({id: n.id, name: n.data.label}));
      if (d.nodes.length > 0) {
        this.setState((state, props) => ({
          namespace: namespaces.find((e) => e.name === this.props.params.namespace),
          namespaces: namespaces
        }));
      }
    });
    this.setState((state, props) => ({
      resource: k8s.resources.find((e) => e.id === this.props.params.kind),
      clusterResource: k8s.cluster.find((e) => e.id === this.props.params.kind)
    }));
  }

  classNames(...classes) {
    return classes.filter(Boolean).join(' ')
  }

  changeNamespace = (namespace) => {
    console.log('nodes changeNamespace: ', namespace);
    //setNamespace(namespace);
    window.open("/namespace/" + namespace.name + "/pod", "_self");
  } 

  changeResource = (resource) => {
    console.log('nodes changeResource: ', resource, this.state.namespace);
    //setResource(resource);
    if (this.state.namespace == undefined) {
      window.open("/namespace/" + this.state.namespaces[0].name + "/" + resource.id, "_self");
    } else {
      window.open("/namespace/" + this.state.namespace.name + "/" + resource.id, "_self"); 
    }
  }

  changeClusterResource = (cr) => {
    console.log('nodes changeClusterResource: ', cr);
    //setClusterResource(cr);
    this.props.onListClick(cr);
  }

  clearSearch = (event) => {
    if (event.target.value === "") {
      this.props.close();
    }
  }

  handleKeyDown = (event) => {
    if (event.key === 'Enter') {
      console.log("event.target.value", event.target.value);
      this.props.onSearchClick({filter: event.target.value});
    }
  }

  signOut = () => {
    fetch('/api/v1/auth/reset', {
      method: 'GET',
    })
    .then(res => res.json())
    .then(d => {
      window.open("/", "_self");
    });
  }

  render() {
    return (
    <Disclosure as="nav" className="bg-white shadow">
      {({ open }) => (
        <>
          <div className="mx-auto max-w-7xl px-2 sm:px-4 lg:px-8">
            <div className="flex h-16 justify-between">
              <div className="flex px-2 lg:px-0">
                <div className="flex flex-shrink-0 items-center">
                  <MenuSelect selected={this.state.namespace || null} setSelected={this.changeNamespace} items={this.state.namespaces} init="Namespace" />
                  <MenuSelectResource selected={this.state.resource || null} setSelected={this.changeResource} items={k8s.resources} init="Resource" />
                  <MenuSelectCluster selected={this.clusterResource || null} setSelected={this.changeClusterResource} items={k8s.cluster} init="Cluster" />
                </div>
              </div>
              <div className="flex flex-1 items-center justify-center px-2 lg:ml-6 lg:justify-end">
                <div className="w-full max-w-lg lg:max-w-xs">
                  <label htmlFor="search" className="sr-only">
                    Search
                  </label>
                  <div className="relative">
                    <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                      <MagnifyingGlassIcon className="h-5 w-5 text-gray-400" aria-hidden="true" />
                    </div>
                    <input
                      id="search"
                      name="search"
                      className="block w-full rounded-md border-0 bg-white py-1.5 pl-10 pr-3 text-gray-900 ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                      placeholder={this.props.filter || "Search"}
                      type="search"
                      onKeyDown={this.handleKeyDown}
                      onChange={this.clearSearch}
                    />
                  </div>
                </div>
              </div>
              <div className="flex items-center lg:hidden">

              </div>
              <div className="hidden lg:ml-4 lg:flex lg:items-center">

              <Menu as="div" className="relative ml-3">
                <div>
                  <Menu.Button className="flex max-w-xs items-center rounded-full bg-white text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
                    <span className="sr-only">Open user menu</span>
                    <div className="h-8 w-8 rounded-full">
                      <UserCircleIcon className="h-8 w-8" aria-hidden="true" />
                    </div>
                  </Menu.Button>
                </div>
                <Transition
                  as={Fragment}
                  enter="transition ease-out duration-200"
                  enterFrom="transform opacity-0 scale-95"
                  enterTo="transform opacity-100 scale-100"
                  leave="transition ease-in duration-75"
                  leaveFrom="transform opacity-100 scale-100"
                  leaveTo="transform opacity-0 scale-95"
                >
                  <Menu.Items className="absolute right-0 z-10 mt-2 w-48 origin-top-right rounded-md bg-white py-1 shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
                    {userNavigation.map((item) => (
                      <Menu.Item key={item.name}>
                        {({ active }) => (
                          <a
                             href="#"
                             onClick={this.signOut}
                             className={this.classNames(
                              active ? 'bg-gray-100' : '',
                              'block px-4 py-2 text-sm text-gray-700'
                            )}
                          >
                            {item.name}
                          </a>
                        )}
                      </Menu.Item>
                    ))}
                  </Menu.Items>
                </Transition>
              </Menu>
              </div>
            </div>
          </div>
        </>
      )}
    </Disclosure>
    );
  }
}

export default Nav;