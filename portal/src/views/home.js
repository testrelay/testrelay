import React from "react";

/* eslint-disable */
const Home = () => (
  <div>
    <div className="container mx-auto p-4">
      <div className="flex justify-between flex- content-center items-center py-2">
        <div className="flex items-center">
          <div className="flex items-center bg-gray-800 shadow w-12 h-12 rounded-full">
            <h1 className="mx-auto font-semibold text-lg text-white">
              <svg height="2rem" viewBox="0 -48 480 480" width="2rem" fill="#fff" xmlns="http://www.w3.org/2000/svg"><path d="m232 0h-64c-3.617188.00390625-6.785156 2.429688-7.726562 5.921875l-45.789063 170.078125h-50.484375c-3.441406 0-6.5 2.203125-7.585938 5.46875l-14.179687 42.53125h-34.234375c-2.121094 0-4.15625.839844-5.65625 2.34375-1.503906 1.5-2.34375 3.535156-2.34375 5.65625v112c0 2.121094.839844 4.15625 2.34375 5.65625 1.5 1.503906 3.535156 2.34375 5.65625 2.34375h69.757812l-5.515624 22.0625c-.597657 2.390625-.0625 4.921875 1.453124 6.859375 1.515626 1.941406 3.84375 3.078125 6.304688 3.078125h64c3.671875 0 6.871094-2.5 7.757812-6.0625l6.484376-25.9375h9.757812c3.8125-.003906 7.09375-2.691406 7.84375-6.429688l14.710938-73.570312h9.445312c2.121094 0 4.15625-.839844 5.65625-2.34375 1.503906-1.5 2.34375-3.535156 2.34375-5.65625v-48c0-2.121094-.839844-4.15625-2.34375-5.65625-1.5-1.503906-3.535156-2.34375-5.65625-2.34375h-13.558594l53.285156-197.921875c.648438-2.402344.140626-4.972656-1.375-6.945313-1.515624-1.976562-3.863281-3.132812-6.351562-3.132812zm-94.25 368h-47.5l4-16h47.5zm54.25-112h-8c-3.8125.003906-7.09375 2.691406-7.84375 6.429688l-14.710938 73.570312h-145.445312v-96h32c3.441406 0 6.5-2.203125 7.585938-5.46875l14.179687-42.53125h40.410156l-5.902343 21.921875c-.648438 2.402344-.140626 4.972656 1.375 6.945313 1.515624 1.976562 3.863281 3.132812 6.351562 3.132812h80zm-22.132812-48h-47.429688l51.695312-192h47.429688zm0 0" /><path d="m472 208h-29.578125l-21.984375-14.65625c-1.3125-.875-2.859375-1.34375-4.4375-1.34375h-168c-2.121094 0-4.15625.839844-5.65625 2.34375-1.503906 1.5-2.34375 3.535156-2.34375 5.65625v128c0 2.121094.839844 4.15625 2.34375 5.65625 1.5 1.503906 3.535156 2.34375 5.65625 2.34375h116.6875l-10.34375 10.34375c-3.125 3.125-3.125 8.1875 0 11.3125l24 24c3.125 3.125 8.1875 3.125 11.3125 0l48-48c.609375-.609375 1.113281-1.308594 1.5-2.078125l5.789062-11.578125h27.054688c2.121094 0 4.15625-.839844 5.65625-2.34375 1.503906-1.5 2.34375-3.535156 2.34375-5.65625v-96c0-2.121094-.839844-4.15625-2.34375-5.65625-1.5-1.503906-3.535156-2.34375-5.65625-2.34375zm-8 96h-24c-3.03125 0-5.800781 1.710938-7.15625 4.421875l-7.421875 14.835937-41.421875 41.429688-12.6875-12.6875 18.34375-18.34375c2.289062-2.289062 2.972656-5.730469 1.734375-8.71875s-4.15625-4.9375-7.390625-4.9375h-128v-16h80c4.417969 0 8-3.582031 8-8s-3.582031-8-8-8h-80v-16h80c4.417969 0 8-3.582031 8-8s-3.582031-8-8-8h-80v-16h80c4.417969 0 8-3.582031 8-8s-3.582031-8-8-8h-80v-16h157.578125l21.984375 14.65625c1.3125.875 2.859375 1.34375 4.4375 1.34375h24zm0 0" /></svg>
            </h1>
          </div>
          <h1 className="ml-2 text-xl font-extrabold">TestRelay</h1>
        </div>

        <div className="hidden md:block">
          <ul className="flex flex-row">
            <li>
              <a className="px-4 py-2 text-gray-700 hover:text-gray-900 transition duration-500 ease-in-out hover:rounded hover:bg-gray-200 rounded-xl" href="#features">Features</a>
            </li>
            <li>
              <a className="px-4 py-2 text-gray-700 hover:text-gray-900 transition duration-500 ease-in-out hover:rounded hover:bg-gray-200 rounded-xl" href="#pricing">Pricing</a>
            </li>
          </ul>
        </div>

        <div className="hidden md:block">
          <ul className="flex flex-">
            <li className="">
              <a className="flex items-center justify-center px-4 py-3 text-base font-medium rounded-md hover:text-gray-500" href="login.html">Login</a>
            </li>
            <li className="px-2">
              <a className="shadow-md flex items-center justify-center px-4 py-3 text-base font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-700" href="signup.html">Signup for free</a>
            </li>
          </ul>
        </div>

        <div className="md:hidden shadow-md flex items-center justify-center px-4 py-3 text-base font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-700">
          Menu
        </div>
        <div className="navbar mobile-nav px-0 mx-0 hidden md:hidden fixed top-0 left-0 w-full bg-white h-screen fixed z-50 p-3">
          <div className="flex flex- justify-between px-3 py-2">
            <img src="img/logo.png" className="w-32 self-start ml-1" />
            <div className="close-menu flex items-center content-center justify-center px-2 py-1 bg-black rounded px-x py-1 text-white uppercase">
              Close
            </div>
          </div>
          <ul className="flex flex-col text-center mt-2 pt-2 w-full">
            <li className="active w-full">
              <a className="w-full font-bold text-lg border-t border-gray-200 block py-3" href="/">Home</a>
            </li>
            <li className="w-full">
              <a className="w-full text-lg border-t border-gray-200 block py-3" href="features.html">Features</a>
            </li>
            <li className="w-full">
              <a className="w-full text-lg border-t border-gray-200 block py-3" href="pricing.html">Pricing</a>
            </li>
            <li className="w-full">
              <a className="w-full text-lg border-t border-gray-200 block py-3" href="contact.html">Contact</a>
            </li>
            <li className="w-full">
              <a className="w-full text-lg border-t border-gray-200 block py-3" href="login.html">Login</a>
            </li>
            <li className="signup py-4 border-t border-gray-200 p-4">
              <a className=" px-3 py-2 bg-yellow-500 rounded border border-yellow-600 shadow font-semibold block" href="signup.html">Signup for free</a>
            </li>

          </ul>
        </div>
      </div>
    </div>

    <div className="mx-auto max-w-7xl sm:px-6 lg:px-8">
      <div className="relative px-4 py-16 overflow-hidden sm:px-0">
        <div>
          <div className="flex flex-col items-center justify-center w-full h-full text-center ">
            <h1 className="mt-4 text-6xl text-gray-800 font-extrabold leading-none tracking-tight">
              Automate Your Take Home Assignments
            </h1>



            <h2 className="mt-3 text-lg text-gray-600 sm:mt-8 sm:max-w-2xl sm:mx-auto md:mt-8 md:text-xl lg:mx-0">
              TestRelay takes care of <span className="text-indigo-500">scheduling, creating</span> and <span className="text-indigo-500">running</span> your take home assignments. Add your candidates and let TestRelay handle the rest.
            </h2>
            <div className="w-full mt-5 sm:mt-12 sm:flex items-center justify-center">
              <a className="shadow-md flex items-center justify-center px-4 py-3 text-base font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-700" href="https://app.slashapi.com/register">
                Get Started Now - It's free
              </a>
            </div>
            <div className="mt-4 w-full">
              <div className="max-w-3xl mx-auto my-10 shadow-md rounded-b-lg">
                <div className="w-full h-11 rounded-t-lg bg-gray-200 flex justify-start items-center space-x-1.5 px-3">
                  <span className="w-3 h-3 rounded-full bg-red-400"></span>
                  <span className="w-3 h-3 rounded-full bg-yellow-400"></span>
                  <span className="w-3 h-3 rounded-full bg-green-400"></span>
                </div>
                <div className="bg-gray-100 border-t-0 w-full  rounded-b-lg flex overflow-hidden" style={{ height: "28rem" }}>
                  <div class="pl-14 text-left">
                    <div className="relative wrap overflow-hidden h-full">
                      <div className="absolute border-gray-200 h-full border-2 ml-4"></div>
                      <div className=" mt-4 flex items-center w-full">
                        <div className="z-20 flex items-center order-1 bg-indigo-500 w-8 h-8 rounded-full">
                          <h1 className="mx-auto font-semibold text-lg text-white">
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                            </svg>
                          </h1>
                        </div>
                        <div className="order-1 bg-white rounded-lg w-3/4 px-6 py-4 ml-8">
                          <h3 className="text-gray-800 text-md capitalize"></h3>
                          <h4 className="mb-3 text-gray-300 text-sm">5 days ago</h4>
                          <p className="text-sm leading-relaxed tracking-wide text-gray-700">TestRelay sent <span className="text-indigo-500">jimmy@code.sh</span> the <span className="font-semibold">BE engineering test</span> with a time limit of 3hrs. They have until the 6th of September to respond.</p>
                        </div>
                      </div>
                      <div className=" mt-4 flex items-center w-full">
                        <div className="z-20 flex items-center order-1 bg-indigo-500 w-8 h-8 rounded-full">
                          <h1 className="mx-auto font-semibold text-lg text-white">
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                            </svg>
                          </h1>
                        </div>
                        <div className="order-1 bg-white rounded-lg w-3/4 px-6 py-4 ml-8">
                          <h3 className="text-gray-800 text-md capitalize"></h3>
                          <h4 className="mb-3 text-gray-300 text-sm">4 days ago</h4>
                          <p className="text-sm leading-relaxed tracking-wide text-gray-700"><span className="text-indigo-500">jimmy@code.sh</span> scheduled a test for 4th September @ 4PM EST.</p>
                        </div>
                      </div>
                      <div className=" mt-4 flex items-center w-full">
                        <div className="z-20 flex items-center order-1 bg-indigo-500 w-8 h-8 rounded-full">
                          <h1 className="mx-auto font-semibold text-lg text-white">
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                            </svg>
                          </h1>
                        </div>
                        <div className="order-1 bg-white rounded-lg w-3/4 px-6 py-4 ml-8">
                          <h3 className="text-gray-800 text-md capitalize"></h3>
                          <h4 className="mb-3 text-gray-300 text-sm">3 days ago @ 16:00</h4>
                          <p className="text-sm leading-relaxed tracking-wide text-gray-700">TestRelay has invited <span className="text-indigo-500">@jimmycode</span> to the assignment repository <span className="text-indigo-500 underline">https://github.com/testrelay/jimmycode-xyx</span>. Test is in progress.</p>
                        </div>
                      </div>
                      <div className=" mt-4 flex items-center w-full">
                        <div className="z-20 flex items-center order-1 bg-green-500 w-8 h-8 rounded-full">
                          <h1 className="mx-auto font-semibold text-lg text-white">
                            <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                            </svg>
                          </h1>
                        </div>
                        <div className="order-1 bg-white rounded-lg w-3/4 px-6 py-4 ml-8">
                          <h3 className="text-gray-800 text-md capitalize"></h3>
                          <h4 className="mb-3 text-gray-300 text-sm">3 days ago @ 20:00</h4>
                          <p className="text-sm leading-relaxed tracking-wide text-gray-700"><span className="text-indigo-500">@jimmycode</span> has submitted his assement to the BE engineering test. Code is ready for review.</p>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div>
      <div className="bg-gray-100 shadow-inner">
        <div className="container max-w-screen-xl mx-auto">
          <div className="py-24 px-8 text-center">
            <h2 className="font-bold text-4xl">Find the best talent with ease</h2>
            <p className="mt-2">1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also </p>
            <div className="">
              <div className="mt-8 font-medium flex justify-center items-center space-x-5 flex-row">
                <div className="flex justify-center items-center">
                  <div className="flex-shrink-0">
                    <svg className="w-5 h-5 text-indigo-600" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
                    </svg>
                  </div>
                  <p className="ml-3 leading-2 text-gray-600">
                    Build functional APIs with zero coding.
                  </p>
                </div>
                <div className="flex justify-center items-center">
                  <div className="flex-shrink-0">
                    <svg className="w-5 h-5 text-indigo-600" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
                    </svg>
                  </div>
                  <p className="ml-3 leading-2 text-gray-600">
                    Resources with permissions.
                  </p>
                </div>
                <div className="flex justify-center items-center">
                  <div className="flex-shrink-0">
                    <svg className="w-5 h-5 text-indigo-600" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
                    </svg>
                  </div>
                  <p className="ml-3 leading-2 text-gray-600">
                    Built in user authentication.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="relative max-w-screen-xl p-4 px-4 py-16 mt-24 mx-auto bg-white dark:bg-gray-800 sm:px-6 lg:px-8">
        <div className="relative">
          <div className="lg:grid lg:grid-flow-row-dense lg:grid-cols-2 lg:gap-8 lg:items-center">
            <div className="ml-auto lg:ml-6 lg:col-start-2 lg:max-w-2xl">
              <p className="text-base font-semibold leading-6 text-indigo-500 uppercase">
                Reduce Development Time
              </p>
              <h4 className="mt-2 text-2xl font-extrabold leading-8 text-gray-900 dark:text-white sm:text-3xl sm:leading-9">
                Evaluate Your Candidates Through Github
              </h4>
              <p className="mt-4 text-lg leading-6 text-gray-600">
                TestRelay hosts and runs code exams on github.
              </p>

              <p className="mt-4 text-lg leading-6 text-gray-600">
                Reduce your development time and increase your speed to market by focus on core app logic not on building APIs.
              </p>

              <p className="mt-4 text-lg leading-6 text-gray-600">
                You take care of the <span className="font-bold text-gray-700">logic</span>, SlashApi will take care of the <span className="font-bold text-gray-700">REST.</span>
              </p>
            </div>
            <div className="relative mt-10 lg:-mx-4 relative-20 lg:mt-0 lg:col-start-1">
              <div className="relative space-y-4">
                <div className="flex items-end justify-center space-x-4 lg:justify-start">
                  <img className="w-32 border rounded-lg shadow-md md:w-56" width="200" src="https://slashapi.com/images/api-token-manager.png" alt="1" />
                  <img className="w-40 rounded-lg shadow-md md:w-64" width="260" src="https://slashapi.com/images/axios-example.png" alt="2" />
                </div>
                <div className="flex items-start justify-center ml-12 space-x-4 lg:justify-start">
                  <img className="w-24 rounded-lg shadow-md md:w-40" width="170" src="https://slashapi.com/images/json-response.png" alt="3" />
                  <img className="w-32 rounded-lg shadow-md md:w-56" width="200" src="https://slashapi.com/images/team-manager.png" alt="4" />
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div className="relative max-w-screen-xl p-4 px-4 py-16 mx-auto bg-white dark:bg-gray-800 sm:px-6 lg:px-8">
        <div className="relative">
          <div className="lg:grid lg:grid-flow-row-dense lg:grid-cols-2 lg:gap-8 lg:items-center">
            <div className="ml-auto lg:mr-6 lg:col-start-1 lg:max-w-2xl">
              <p className="text-base font-semibold leading-6 text-indigo-500 uppercase">
                Build APIs in minutes
              </p>
              <h4 className="mt-2 text-2xl font-extrabold leading-8 text-gray-900 dark:text-white sm:text-3xl sm:leading-9">
                Instant APIs without code.
              </h4>
              <p className="mt-4 text-lg leading-6 text-gray-600">
                SlashApi instantly build a secure, reusable, fully documented and functional APIs with zero coding.
              </p>

              <p className="mt-4 text-lg leading-6 text-gray-600">
                Build a REST API from your favorite apps and tools with a few clicks. Everyone can build APIs, no technical or backend coding skills required.
              </p>
            </div>
            <div className="relative mt-10 lg:-mx-4 relative-20 lg:mt-0 lg:col-start-2">
              <img src="https://slashapi.com/images/instant-api.png" alt="" />
            </div>
          </div>
        </div>
      </div>

      <div className="relative max-w-screen-xl p-4 px-4 py-16 mx-auto bg-white dark:bg-gray-800 sm:px-6 lg:px-8">
        <div className="relative">
          <div className="lg:grid lg:grid-flow-row-dense lg:grid-cols-2 lg:gap-8 lg:items-center">
            <div className="ml-auto lg:ml-6 lg:col-start-2 lg:max-w-2xl">
              <p className="text-base font-semibold leading-6 text-indigo-500 uppercase">
                Securing your API
              </p>
              <h4 className="mt-2 text-2xl font-extrabold leading-8 text-gray-900 dark:text-white sm:text-3xl sm:leading-9">
                Powerful Security Options
              </h4>

              <ul className="mt-8 space-y-3 font-medium">
                <li className="flex items-start lg:col-span-1">
                  <div className="flex-shrink-0">
                    <svg className="w-5 h-5 text-indigo-600" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
                    </svg>
                  </div>
                  <p className="ml-3 leading-5 text-gray-600">
                    Secure every API endpoint with our built in authentication.
                  </p>
                </li>
                <li className="flex items-start lg:col-span-1">
                  <div className="flex-shrink-0">
                    <svg className="w-5 h-5 text-indigo-600" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
                    </svg>
                  </div>
                  <p className="ml-3 leading-5 text-gray-600">
                    Restrict capabilities using Token-Based Access Controls.
                  </p>
                </li>
                <li className="flex items-start lg:col-span-1">
                  <div className="flex-shrink-0">
                    <svg className="w-5 h-5 text-indigo-600" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
                    </svg>
                  </div>
                  <p className="ml-3 leading-5 text-gray-600">
                    Easily Manage API Keys.
                  </p>
                </li>
                <li className="flex items-start lg:col-span-1">
                  <div className="flex-shrink-0">
                    <svg className="w-5 h-5 text-indigo-600" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
                    </svg>
                  </div>
                  <p className="ml-3 leading-5 text-gray-600">
                    Custom validations to validate input coming to your API.
                  </p>
                </li>
              </ul>
            </div>
            <div className="relative mt-10 relative-20 lg:mt-0 lg:col-start-1">
              <img src="https://slashapi.com/images/security-options.png" alt="Security Options" className="p-2 border" />
            </div>
          </div>
        </div>
      </div>

      <div className="relative max-w-screen-xl p-4 px-4 py-16 mx-auto bg-white dark:bg-gray-800 sm:px-6 lg:px-8">
        <div className="relative">
          <div className="lg:grid lg:grid-flow-row-dense lg:grid-cols-2 lg:gap-8 lg:items-center">
            <div className="ml-auto lg:mr-6 lg:max-w-2xl">
              <p className="text-base font-semibold leading-6 text-indigo-500 uppercase">
                Integrations
              </p>
              <h4 className="mt-2 text-2xl font-extrabold leading-8 text-gray-900 dark:text-white sm:text-3xl sm:leading-9">
                Integrated with common apps
              </h4>
              <p className="mt-4 text-lg leading-6 text-gray-600">
                Connect SlashApi with the apps and tools you use every day.
              </p>
              <p className="mt-4 text-lg leading-6 text-gray-600">

                Choose from our app collections to work with all of your data using REST APIs.
              </p>

              <div className="mt-8">
                <a href="https://slashapi.com/collections" className="inline-flex items-center px-4 py-2 text-base font-medium text-white bg-indigo-600 border border-transparent rounded-md shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                  Explore All
                  <svg className="w-5 h-5 ml-3 -mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path></svg>
                </a>
              </div>
            </div>
            <div className="relative mt-10 space-y-6 relative-20 lg:mt-0">
              <div className="grid grid-cols-5">
                <div>&nbsp;</div>
                <div className="flex items-center justify-center w-16 h-16 bg-gray-100 lg:w-20 lg:h-20 rounded-2xl">
                  <div className="w-12 h-12 lg:w-16 lg:h-16">
                    <img className="object-contain h-full" src="https://slashapi.com/images/services/mysql.svg" alt="MySQL" />
                  </div>
                </div>
                <div className="flex items-center justify-center w-16 h-16 lg:w-20 lg:h-20 rounded-2xl" >
                  <div className="w-12 h-12 lg:w-16 lg:h-16">
                    <img className="object-contain h-full" src="https://slashapi.com/images/services/pgsql.svg" alt="Pgsql" />
                  </div>
                </div>
                <div className="flex items-center justify-center w-16 h-16 bg-gray-100 lg:w-20 lg:h-20 rounded-2xl">
                  <div className="w-12 h-12 lg:w-16 lg:h-16">
                    <img className="relative object-contain h-full" src="https://slashapi.com/images/services/mongodb.svg" alt="MongoDB" />
                  </div>
                </div>
                <div>&nbsp;</div>
              </div>

              <div className="grid grid-cols-5">
                <div className="flex items-center justify-center w-16 h-16 bg-gray-100 lg:w-20 lg:h-20 rounded-2xl">
                  <div className="w-12 h-12 lg:w-14 lg:h-14">
                    <img className="object-contain h-full" src="https://slashapi.com/images/services/azure.svg" alt="Azure" />
                  </div>
                </div>
                <div className="flex items-center justify-center w-16 h-16 bg-gray-100 lg:w-20 lg:h-20 rounded-2xl">
                  <div className="w-12 h-12 lg:w-16 lg:h-16">
                    <img className="object-contain h-full" src="https://slashapi.com/images/services/dropbox.svg" alt="Dropbox" />
                  </div>
                </div>
                <div className="flex items-center justify-center w-16 h-16 bg-gray-100 lg:w-20 lg:h-20 rounded-2xl">
                  <div className="w-12 h-12 lg:w-16 lg:h-16">
                    <img className="object-contain h-full" src="https://slashapi.com/images/services/ftp.svg" alt="FTP" />
                  </div>
                </div>
                <div className="flex items-center justify-center w-16 h-16 bg-gray-100 lg:w-20 lg:h-20 rounded-2xl">
                  <div className="w-12 h-12 lg:w-16 lg:h-16">
                    <img className="object-contain h-full" src="https://slashapi.com/images/services/aws-s3.svg" alt="AWS S3" />
                  </div>
                </div>
                <div className="flex items-center justify-center w-16 h-16 bg-gray-100 lg:w-20 lg:h-20 rounded-2xl">
                  <div className="w-12 h-12 lg:w-14 lg:h-14">
                    <img className="object-contain h-full" src="https://slashapi.com/images/services/spaces.svg" alt="DO Spaces" />
                  </div>
                </div>
              </div>

              <div className="grid grid-cols-5">
                <div className="flex items-center justify-center w-16 h-16 bg-gray-100 lg:w-20 lg:h-20 rounded-2xl">
                  <div className="w-12 h-12 lg:w-14 lg:h-14">
                    <img className="object-contain h-full" src="https://slashapi.com/images/services/twitter.svg" alt="Twitter" />
                  </div>
                </div>
                <div className="flex items-center justify-center w-16 h-16 bg-gray-100 lg:w-20 lg:h-20 rounded-2xl">
                  <div className="w-12 h-12 lg:w-14 lg:h-14">
                    <img className="object-contain h-full" src="https://slashapi.com/images/services/google-calendar.svg" alt="Google Calendar" />
                  </div>
                </div>
                <div className="flex items-center justify-center w-16 h-16 bg-gray-100 lg:w-20 lg:h-20 rounded-2xl">
                  <div className="w-12 h-12 lg:w-16 lg:h-16">
                    <img className="object-contain h-full" src="https://slashapi.com/images/services/telegram.svg" alt="Telegram" />
                  </div>
                </div>
                <div className="flex items-center justify-center w-16 h-16 bg-gray-100 lg:w-20 lg:h-20 rounded-2xl">
                  <div className="w-12 h-12 lg:w-14 lg:h-14">
                    <img className="object-contain h-full" src="https://slashapi.com/images/services/slack.svg" alt="Slack" />
                  </div>
                </div>
                <div className="flex items-center justify-center w-16 h-16 bg-gray-100 lg:w-20 lg:h-20 rounded-2xl">
                  <div className="w-12 h-12 lg:w-14 lg:h-14">
                    <img className="relative object-contain h-full left-2" src="https://slashapi.com/images/services/google-sheets.svg" alt="google-sheets" />
                  </div>
                </div>
              </div>

              <div className="grid grid-cols-5">
                <div>&nbsp;</div>
                <div className="flex items-center justify-center w-16 h-16 bg-gray-100 lg:w-20 lg:h-20 rounded-2xl">
                  <div className="w-12 h-12 lg:w-16 lg:h-16">
                    <img className="object-contain h-full" src="https://slashapi.com/images/services/twilio.svg" alt="Twilio" />
                  </div>
                </div>
                <div className="flex items-center justify-center w-16 h-16 bg-gray-100 lg:w-20 lg:h-20 rounded-2xl">
                  <div className="w-12 h-12 lg:w-16 lg:h-16">
                    <img className="relative object-contain h-full left-1" src="https://slashapi.com/images/services/xml-to-json.svg" alt="XML" />
                  </div>
                </div>
                <div className="flex items-center justify-center w-16 h-16 bg-gray-100 lg:w-20 lg:h-20 rounded-2xl">
                  <div className="w-10 h-10 lg:w-14 lg:h-14">
                    <img className="object-contain h-full" src="https://slashapi.com/images/services/airtable.svg" alt="airtable" />
                  </div>
                </div>
                <div>&nbsp;</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div className="mt-24 bg-gray-100 dark:bg-gray-800 ">
      <div className="relative max-w-screen-xl p-4 px-4 py-16 mx-auto  sm:px-6 lg:px-8">
        <section className="text-gray-600 body-font overflow-hidden">
          <div className="container px-5 py-8 mx-auto">
            <div className="flex flex-col text-center w-full mb-12">
              <h2 className="sm:text-4xl text-3xl font-medium text-bold mb-2 text-gray-800">Pricing</h2>
              <p className="lg:w-2/3 mx-auto leading-relaxed text-base text-gray-500">We're committed to having a free plan, forever. Only pay for features when you can afford to do so.</p>
            </div>
            <div className="flex flex-wrap justify-center -m-4">
              <div className="p-4 xl:w-1/4 md:w-1/2 w-full">
                <div className="h-full p-6 rounded-lg bg-white  border-2 border-indigo-500 shadow-md flex flex-col relative overflow-hidden">
                  <h2 className="text-sm tracking-widest title-font mb-1 font-medium">START</h2>
                  <h1 className="text-5xl text-gray-900 pb-4 mb-4 border-b border-gray-200 leading-none">Free</h1>
                  <p className="flex items-center text-gray-600 mb-2">
                    <span className="w-4 h-4 mr-2 inline-flex items-center justify-center bg-indigo-600 text-white rounded-full flex-shrink-0">
                      <svg fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" className="w-3 h-3" viewBox="0 0 24 24">
                        <path d="M20 6L9 17l-5-5"></path>
                      </svg>
                    </span>Unlimited assignments
                  </p>
                  <p className="flex items-center text-gray-600 mb-2">
                    <span className="w-4 h-4 mr-2 inline-flex items-center justify-center bg-indigo-600 text-white rounded-full flex-shrink-0">
                      <svg fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" className="w-3 h-3" viewBox="0 0 24 24">
                        <path d="M20 6L9 17l-5-5"></path>
                      </svg>
                    </span>Unlimited users
                  </p>
                  <p className="flex items-center text-gray-600 mb-6">
                    <span className="w-4 h-4 mr-2 inline-flex items-center justify-center bg-indigo-600 text-white rounded-full flex-shrink-0">
                      <svg fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" className="w-3 h-3" viewBox="0 0 24 24">
                        <path d="M20 6L9 17l-5-5"></path>
                      </svg>
                    </span>Everything for free!
                  </p>
                  <button className="flex items-center mt-auto text-white bg-indigo-500 bborder-0 py-2 px-4 w-full focus:outline-none hover:bg-gray-500 rounded">Get Started
                    <svg fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" className="w-4 h-4 ml-auto" viewBox="0 0 24 24">
                      <path d="M5 12h14M12 5l7 7-7 7"></path>
                    </svg>
                  </button>
                  <p className="text-xs text-gray-500 mt-3">What are you waiting for? It's free!</p>
                </div>
              </div>
              <div className="p-4 xl:w-1/4 md:w-1/2 w-full">
                <div className="h-full p-6 rounded-lg  bg-white shadow-md flex flex-col relative overflow-hidden">
                  <span className="bg-gray-300 text-white px-3 py-1 tracking-widest text-xs absolute right-0 top-0 rounded-bl">Coming Soon</span>
                  <h2 className="text-sm tracking-widest title-font mb-1 font-medium">PRO</h2>
                  <h1 className="text-5xl text-gray-900 leading-none flex items-center pb-4 mb-4 border-b border-gray-200">
                    <span >$?</span>
                    <span className="text-lg ml-1 font-normal text-gray-500">/mo</span>
                  </h1>
                  <p className="flex items-center text-gray-600 mb-2">
                    We're focused on building a product that people will love rather than figuring out pricing. All plans are free for now and we will always have a free tier.
                  </p>
                  <p className="text-xs text-gray-500 mt-3">Want a say in what's in the pro plan, email <span className="text-indigo-500">hello@testrelay.io.</span></p>
                </div>
              </div>
            </div>
          </div>
        </section>
      </div>
    </div>
    <footer className="text-gray-600 body-font">
      <div className="container px-5 py-8 mx-auto flex items-center sm:flex-row flex-col">
        <a className="flex title-font font-medium items-center md:justify-start justify-center text-gray-900">
          <div className="flex items-center bg-gray-800 shadow w-10 h-10 rounded-full">
            <h1 className="mx-auto font-semibold text-lg text-white">
              <svg height="1.8rem" viewBox="0 -48 480 480" width="1.8rem" fill="#fff" xmlns="http://www.w3.org/2000/svg"><path d="m232 0h-64c-3.617188.00390625-6.785156 2.429688-7.726562 5.921875l-45.789063 170.078125h-50.484375c-3.441406 0-6.5 2.203125-7.585938 5.46875l-14.179687 42.53125h-34.234375c-2.121094 0-4.15625.839844-5.65625 2.34375-1.503906 1.5-2.34375 3.535156-2.34375 5.65625v112c0 2.121094.839844 4.15625 2.34375 5.65625 1.5 1.503906 3.535156 2.34375 5.65625 2.34375h69.757812l-5.515624 22.0625c-.597657 2.390625-.0625 4.921875 1.453124 6.859375 1.515626 1.941406 3.84375 3.078125 6.304688 3.078125h64c3.671875 0 6.871094-2.5 7.757812-6.0625l6.484376-25.9375h9.757812c3.8125-.003906 7.09375-2.691406 7.84375-6.429688l14.710938-73.570312h9.445312c2.121094 0 4.15625-.839844 5.65625-2.34375 1.503906-1.5 2.34375-3.535156 2.34375-5.65625v-48c0-2.121094-.839844-4.15625-2.34375-5.65625-1.5-1.503906-3.535156-2.34375-5.65625-2.34375h-13.558594l53.285156-197.921875c.648438-2.402344.140626-4.972656-1.375-6.945313-1.515624-1.976562-3.863281-3.132812-6.351562-3.132812zm-94.25 368h-47.5l4-16h47.5zm54.25-112h-8c-3.8125.003906-7.09375 2.691406-7.84375 6.429688l-14.710938 73.570312h-145.445312v-96h32c3.441406 0 6.5-2.203125 7.585938-5.46875l14.179687-42.53125h40.410156l-5.902343 21.921875c-.648438 2.402344-.140626 4.972656 1.375 6.945313 1.515624 1.976562 3.863281 3.132812 6.351562 3.132812h80zm-22.132812-48h-47.429688l51.695312-192h47.429688zm0 0" /><path d="m472 208h-29.578125l-21.984375-14.65625c-1.3125-.875-2.859375-1.34375-4.4375-1.34375h-168c-2.121094 0-4.15625.839844-5.65625 2.34375-1.503906 1.5-2.34375 3.535156-2.34375 5.65625v128c0 2.121094.839844 4.15625 2.34375 5.65625 1.5 1.503906 3.535156 2.34375 5.65625 2.34375h116.6875l-10.34375 10.34375c-3.125 3.125-3.125 8.1875 0 11.3125l24 24c3.125 3.125 8.1875 3.125 11.3125 0l48-48c.609375-.609375 1.113281-1.308594 1.5-2.078125l5.789062-11.578125h27.054688c2.121094 0 4.15625-.839844 5.65625-2.34375 1.503906-1.5 2.34375-3.535156 2.34375-5.65625v-96c0-2.121094-.839844-4.15625-2.34375-5.65625-1.5-1.503906-3.535156-2.34375-5.65625-2.34375zm-8 96h-24c-3.03125 0-5.800781 1.710938-7.15625 4.421875l-7.421875 14.835937-41.421875 41.429688-12.6875-12.6875 18.34375-18.34375c2.289062-2.289062 2.972656-5.730469 1.734375-8.71875s-4.15625-4.9375-7.390625-4.9375h-128v-16h80c4.417969 0 8-3.582031 8-8s-3.582031-8-8-8h-80v-16h80c4.417969 0 8-3.582031 8-8s-3.582031-8-8-8h-80v-16h80c4.417969 0 8-3.582031 8-8s-3.582031-8-8-8h-80v-16h157.578125l21.984375 14.65625c1.3125.875 2.859375 1.34375 4.4375 1.34375h24zm0 0" /></svg>
            </h1>
          </div>
          <span className="ml-3 text-xl">TestRelay</span>
        </a>
        <p className="text-sm text-gray-500 sm:ml-4 sm:pl-4 sm:border-l-2 sm:border-gray-200 sm:py-2 sm:mt-0 mt-4">© 2021 TestRelay —
          <a href="https://twitter.com/knyttneve" className="text-gray-600 ml-1" rel="noopener noreferrer" target="_blank">@hugorut</a>
        </p>
        <span className="inline-flex sm:ml-auto sm:mt-0 mt-4 justify-center sm:justify-start">
          <a className="text-gray-500">
            <svg fill="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" className="w-5 h-5" viewBox="0 0 24 24">
              <path d="M18 2h-3a5 5 0 00-5 5v3H7v4h3v8h4v-8h3l1-4h-4V7a1 1 0 011-1h3z"></path>
            </svg>
          </a>
          <a className="ml-3 text-gray-500">
            <svg fill="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" className="w-5 h-5" viewBox="0 0 24 24">
              <path d="M23 3a10.9 10.9 0 01-3.14 1.53 4.48 4.48 0 00-7.86 3v1A10.66 10.66 0 013 4s-4 9 5 13a11.64 11.64 0 01-7 2c9 5 20 0 20-11.5a4.5 4.5 0 00-.08-.83A7.72 7.72 0 0023 3z"></path>
            </svg>
          </a>
          <a className="ml-3 text-gray-500">
            <svg fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" className="w-5 h-5" viewBox="0 0 24 24">
              <rect width="20" height="20" x="2" y="2" rx="5" ry="5"></rect>
              <path d="M16 11.37A4 4 0 1112.63 8 4 4 0 0116 11.37zm1.5-4.87h.01"></path>
            </svg>
          </a>
          <a className="ml-3 text-gray-500">
            <svg fill="currentColor" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="0" className="w-5 h-5" viewBox="0 0 24 24">
              <path stroke="none" d="M16 8a6 6 0 016 6v7h-4v-7a2 2 0 00-2-2 2 2 0 00-2 2v7h-4v-7a6 6 0 016-6zM2 9h4v12H2z"></path>
              <circle cx="4" cy="4" r="2" stroke="none"></circle>
            </svg>
          </a>
        </span>
      </div>
    </footer>
  </div >
);

export default Home;
