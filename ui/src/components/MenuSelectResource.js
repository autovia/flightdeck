// Copyright (c) Autovia GmbH
// SPDX-License-Identifier: Apache-2.0

import { Fragment, useState, useEffect } from 'react'
import { Listbox, Transition } from '@headlessui/react'
import { CheckIcon, ChevronDownIcon } from '@heroicons/react/20/solid'

function classNames(...classes) {
  return classes.filter(Boolean).join(' ')
}

export default function MenuSelect({selected, setSelected, items, init}) {
  useEffect(() => {
    console.log('selected', selected);
  }, [selected]);

  return (
    <Listbox value={selected} onChange={setSelected}>
      {({ open }) => (
        <>
          <div className="relative">
            <Listbox.Button className="relative w-full cursor-default rounded-md bg-white py-1.5 pl-3 pr-10 text-left text-gray-900 shadow-base base:text-base sm:leading-6">
              <span className={selected ? "block truncate text-base font-semibold text-gray-900 shadow-base" : "block truncate"}
              >{selected ? selected.name : init}</span>
              <span className="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
              <ChevronDownIcon className="-mr-1 h-5 w-5 text-gray-400" aria-hidden="true" />
              </span>
            </Listbox.Button>

            <Transition
              show={open}
              as={Fragment}
              leave="transition ease-in duration-100"
              leaveFrom="opacity-100"
              leaveTo="opacity-0"
            >
              <Listbox.Options className="absolute left-0 z-10 mt-2 w-300 origin-top-left rounded-md bg-white shadow-lg ring-1 ring-black ring-opacity-5 focus:outline-none">
              <div className="table w-full">
                <div className="table-header-group">
                  <div className="table-row">
                    <div className="table-cell text-left p-3 font-semibold">Workloads</div>
                    <div className="table-cell text-left p-3 font-semibold">Services</div>
                    <div className="table-cell text-left p-3 font-semibold">Cluster</div>
                  </div>
                </div>
                <div className="table-row-group">
                  <div className="table-row">
                    <div className="table-cell">

                    {items.filter((i) => i.type==="workload").map((item) => (
                      <Listbox.Option
                        key={item.id}
                        className={({ active }) =>
                          classNames(
                            active ? 'bg-sky-600 text-white' : 'text-gray-900',
                            'relative cursor-default select-none py-2 pl-3 pr-9'
                          )
                        }
                        value={item}
                      >
                        {({ selected, active }) => (
                          <>
                            <span className={classNames(selected ? 'font-semibold' : 'font-normal', 'block truncate')}>
                              {item.name}
                            </span>

                            {selected ? (
                              <span
                                className={classNames(
                                  active ? 'text-white' : 'text-sky-600',
                                  'absolute inset-y-0 right-0 flex items-center pr-4'
                                )}
                              >
                                <CheckIcon className="h-5 w-5" aria-hidden="true" />
                              </span>
                            ) : null}
                          </>
                        )}
                      </Listbox.Option>
                    ))}

                    </div>
                    <div className="table-cell">

                    {items.filter((i) => i.type==="service").map((item) => (
                      <Listbox.Option
                        key={item.id}
                        className={({ active }) =>
                          classNames(
                            active ? 'bg-sky-600 text-white' : 'text-gray-900',
                            'relative cursor-default select-none py-2 pl-3 pr-9'
                          )
                        }
                        value={item}
                      >
                        {({ selected, active }) => (
                          <>
                            <span className={classNames(selected ? 'font-semibold' : 'font-normal', 'block truncate')}>
                              {item.name}
                            </span>

                            {selected ? (
                              <span
                                className={classNames(
                                  active ? 'text-white' : 'text-sky-600',
                                  'absolute inset-y-0 right-0 flex items-center pr-4'
                                )}
                              >
                                <CheckIcon className="h-5 w-5" aria-hidden="true" />
                              </span>
                            ) : null}
                          </>
                        )}
                      </Listbox.Option>
                    ))}

                    <div className="table-cell text-left p-3 pt-12 font-semibold">Configuration</div>

                    {items.filter((i) => i.type==="config").map((item) => (
                      <Listbox.Option
                        key={item.id}
                        className={({ active }) =>
                          classNames(
                            active ? 'bg-sky-600 text-white' : 'text-gray-900',
                            'relative cursor-default select-none py-2 pl-3 pr-9'
                          )
                        }
                        value={item}
                      >
                        {({ selected, active }) => (
                          <>
                            <span className={classNames(selected ? 'font-semibold' : 'font-normal', 'block truncate')}>
                              {item.name}
                            </span>

                            {selected ? (
                              <span
                                className={classNames(
                                  active ? 'text-white' : 'text-sky-600',
                                  'absolute inset-y-0 right-0 flex items-center pr-4'
                                )}
                              >
                                <CheckIcon className="h-5 w-5" aria-hidden="true" />
                              </span>
                            ) : null}
                          </>
                        )}
                      </Listbox.Option>
                    ))}

                    </div>
                    <div className="table-cell">

                    {items.filter((i) => i.type==="cluster").map((item) => (
                      <Listbox.Option
                        key={item.id}
                        className={({ active }) =>
                          classNames(
                            active ? 'bg-sky-600 text-white' : 'text-gray-900',
                            'relative cursor-default select-none py-2 pl-3 pr-9'
                          )
                        }
                        value={item}
                      >
                        {({ selected, active }) => (
                          <>
                            <span className={classNames(selected ? 'font-semibold' : 'font-normal', 'block truncate')}>
                              {item.name}
                            </span>

                            {selected ? (
                              <span
                                className={classNames(
                                  active ? 'text-white' : 'text-sky-600',
                                  'absolute inset-y-0 right-0 flex items-center pr-4'
                                )}
                              >
                                <CheckIcon className="h-5 w-5" aria-hidden="true" />
                              </span>
                            ) : null}
                          </>
                        )}
                      </Listbox.Option>
                    ))}

                      <div className="table-cell text-left p-3 pt-12 font-semibold">Storage</div>

                      {items.filter((i) => i.type==="storage").map((item) => (
                        <Listbox.Option
                          key={item.id}
                          className={({ active }) =>
                            classNames(
                              active ? 'bg-sky-600 text-white' : 'text-gray-900',
                              'relative cursor-default select-none py-2 pl-3 pr-9'
                            )
                          }
                          value={item}
                        >
                          {({ selected, active }) => (
                            <>
                              <span className={classNames(selected ? 'font-semibold' : 'font-normal', 'block truncate')}>
                                {item.name}
                              </span>

                              {selected ? (
                                <span
                                  className={classNames(
                                    active ? 'text-white' : 'text-sky-600',
                                    'absolute inset-y-0 right-0 flex items-center pr-4'
                                  )}
                                >
                                  <CheckIcon className="h-5 w-5" aria-hidden="true" />
                                </span>
                              ) : null}
                            </>
                          )}
                        </Listbox.Option>
                      ))}


                    </div>
                  </div>
                </div>
              </div>
              
              </Listbox.Options>              
            </Transition>
          </div>
        </>
      )}
    </Listbox>
  )
}