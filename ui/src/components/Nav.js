// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import { Fragment, useState, useEffect } from 'react'
import { Disclosure, Menu, Transition } from '@headlessui/react'
import { UserCircleIcon } from '@heroicons/react/20/solid'
import { Bars3Icon, BellIcon, XMarkIcon } from '@heroicons/react/24/outline'
import MenuSelect from './MenuSelect';
import MenuSelectResource from './MenuSelectResource';
import MenuSelectCluster from './MenuSelectCluster';

function classNames(...classes) {
  return classes.filter(Boolean).join(' ')
}

const userNavigation = [
  { name: 'Sign out', href: '#' },
]

const resources = [
  { id: "cm", name: 'Config Maps', type: 'config' },
  { id: "cronjob", name: 'Cron Jobs', type: 'workload' },
  { id: "ds", name: 'DaemonSet', type: 'workload' },
  { id: "deploy", name: 'Deployments', type: 'workload' },
  { id: "ev", name: 'Events', type: 'cluster' },
  { id: "ing", name: 'Ingresses', type: 'service' },
  { id: "job", name: 'Jobs', type: 'workload' },
  { id: "netpol", name: 'Network Policies', type: 'cluster' },
  { id: "pvc", name: 'Persistent Volume Claims', type: 'storage' },
  { id: "pod", name: 'Pods', type: 'workload'}, // default
  { id: "rs", name: 'Replica Sets', type: 'workload' },
  { id: "rc", name: 'Replication Controllers', type: 'workload' },
  { id: "role", name: 'Roles', type: 'cluster' },
  { id: "rb", name: 'Role Bindings', type: 'cluster' },
  { id: "secret", name: 'Secrets', type: 'config' },
  { id: "svc", name: 'Services', type: 'service' },
  { id: "sa", name: 'Service Accounts', type: 'cluster' },
  { id: "sts", name: 'Stateful Sets', type: 'workload' }
]

const clusterResources = [
  { id: "c-role", name: 'Cluster Roles', type: 'cluster' },
  { id: "crb", name: 'Cluster Role Bindings', type: 'cluster' },
  { id: "crd", name: 'Custom Resource Definitions', type: 'cluster' },
  { id: "ic", name: 'Ingress Classes', type: 'service' },
  { id: "no", name: 'Nodes', type: 'cluster' }, // default
  { id: "pv", name: 'Persistent Volumes', type: 'storage' },
  { id: "sc", name: 'Storage Classes', type: 'storage' }
]

export default function Nav({params}) {
  const [namespace, setNamespace] = useState({});
  const [namespaces, setNamespaces] = useState([]);
  const [resource, setResource] = useState({});
  const [clusterResource, setClusterResource] = useState({});

  useEffect(() => {
    console.log("params", params);
    fetch('/api/v1/namespaces')
    .then(res => res.json())
    .then(d => {
      console.log('/api/v1/namespaces', d);
      const namespaces = d.nodes.map((n) => ({id: n.id, name: n.data.label}));
      if (d.nodes.length > 0) {
        setNamespaces(namespaces);
        setNamespace(namespaces.find((e) => e.name === params.namespace));
      }
    });
    setResource(resources.find((e) => e.id === params.kind));
    setClusterResource(clusterResources.find((e) => e.id === params.kind));
  }, []);

  const changeNamespace = (namespace) => {
    console.log('nodes changeNamespace: ', namespace);
    //setNamespace(namespace);
    window.open("/namespace/" + namespace.name + "/pod", "_self");
  } 

  const changeResource = (resource) => {
    console.log('nodes changeResource: ', resource, namespace);
    //setResource(resource);
    if (namespace == undefined) {
      window.open("/namespace/" + namespaces[0].name + "/" + resource.id, "_self");
    } else {
      window.open("/namespace/" + namespace.name + "/" + resource.id, "_self"); 
    }
  }

  const changeClusterResource = (cr) => {
    console.log('nodes changeClusterResource: ', cr);
    //setClusterResource(cr);
    window.open("/cluster/" + cr.id, "_self");
  }

  const signOut = () => {
    fetch('/api/v1/auth/reset', {
      method: 'GET',
    })
    .then(res => res.json())
    .then(d => {
      window.open("/", "_self");
    });
  }

  return (
    <Disclosure as="nav" className="bg-white shadow">
      {({ open }) => (
        <>
          <div className="mx-auto max-w-7xl px-2 sm:px-4 lg:px-8">
            <div className="flex h-16 justify-between">
              <div className="flex px-2 lg:px-0">
                <div className="flex flex-shrink-0 items-center">
                  <MenuSelect selected={namespace || null} setSelected={changeNamespace} items={namespaces} init="Namespace" />
                  <MenuSelectResource selected={resource || null} setSelected={changeResource} items={resources} init="Resource" />
                  <MenuSelectCluster selected={clusterResource || null} setSelected={changeClusterResource} items={clusterResources} init="Cluster" />
                </div>
              </div>
              <div className="flex flex-1 items-center justify-center px-2 lg:ml-6 lg:justify-end">

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
                             onClick={signOut}
                             className={classNames(
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
  )
}