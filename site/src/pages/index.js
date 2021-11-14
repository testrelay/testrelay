import * as React from "react"
import {Helmet} from 'react-helmet';
import Changelog from '../images/changelog.png'
import Review from '../images/review.png'
import Cursor from '../images/cursor1.webp'
import DatePicker from '../images/datepicker.png'
import Header from '../images/header.png'

const IndexPage = () => {
    return (
        <div>
            <Helmet>
                <title>TestRelay</title>
                <meta name="description"
                      content="an all-in-one platform for managing your technical interviews through Github"/>
                <meta name="og:image" content={Header}/>
            </Helmet>
            <main>
                <div className="px-4 py-2 bg-blue-600 text-center shadow">
                    <span className="text-white mx-auto">We're in <span className="text-yellow-500">alpha</span>, signup for early access now</span>
                </div>
                <div className="container mx-auto p-4">
                    <div className="flex justify-between flex- content-center items-center py-2">
                        <div className="flex items-center">
                            <div className="flex items-center bg-gray-800 shadow w-12 h-12 rounded-full">
                                <h1 className="mx-auto font-semibold text-lg text-white">
                                    <svg height="2rem" viewBox="0 -48 480 480" width="2rem" fill="#fff"
                                         xmlns="http://www.w3.org/2000/svg">
                                        <path
                                            d="m232 0h-64c-3.617188.00390625-6.785156 2.429688-7.726562 5.921875l-45.789063 170.078125h-50.484375c-3.441406 0-6.5 2.203125-7.585938 5.46875l-14.179687 42.53125h-34.234375c-2.121094 0-4.15625.839844-5.65625 2.34375-1.503906 1.5-2.34375 3.535156-2.34375 5.65625v112c0 2.121094.839844 4.15625 2.34375 5.65625 1.5 1.503906 3.535156 2.34375 5.65625 2.34375h69.757812l-5.515624 22.0625c-.597657 2.390625-.0625 4.921875 1.453124 6.859375 1.515626 1.941406 3.84375 3.078125 6.304688 3.078125h64c3.671875 0 6.871094-2.5 7.757812-6.0625l6.484376-25.9375h9.757812c3.8125-.003906 7.09375-2.691406 7.84375-6.429688l14.710938-73.570312h9.445312c2.121094 0 4.15625-.839844 5.65625-2.34375 1.503906-1.5 2.34375-3.535156 2.34375-5.65625v-48c0-2.121094-.839844-4.15625-2.34375-5.65625-1.5-1.503906-3.535156-2.34375-5.65625-2.34375h-13.558594l53.285156-197.921875c.648438-2.402344.140626-4.972656-1.375-6.945313-1.515624-1.976562-3.863281-3.132812-6.351562-3.132812zm-94.25 368h-47.5l4-16h47.5zm54.25-112h-8c-3.8125.003906-7.09375 2.691406-7.84375 6.429688l-14.710938 73.570312h-145.445312v-96h32c3.441406 0 6.5-2.203125 7.585938-5.46875l14.179687-42.53125h40.410156l-5.902343 21.921875c-.648438 2.402344-.140626 4.972656 1.375 6.945313 1.515624 1.976562 3.863281 3.132812 6.351562 3.132812h80zm-22.132812-48h-47.429688l51.695312-192h47.429688zm0 0"/>
                                        <path
                                            d="m472 208h-29.578125l-21.984375-14.65625c-1.3125-.875-2.859375-1.34375-4.4375-1.34375h-168c-2.121094 0-4.15625.839844-5.65625 2.34375-1.503906 1.5-2.34375 3.535156-2.34375 5.65625v128c0 2.121094.839844 4.15625 2.34375 5.65625 1.5 1.503906 3.535156 2.34375 5.65625 2.34375h116.6875l-10.34375 10.34375c-3.125 3.125-3.125 8.1875 0 11.3125l24 24c3.125 3.125 8.1875 3.125 11.3125 0l48-48c.609375-.609375 1.113281-1.308594 1.5-2.078125l5.789062-11.578125h27.054688c2.121094 0 4.15625-.839844 5.65625-2.34375 1.503906-1.5 2.34375-3.535156 2.34375-5.65625v-96c0-2.121094-.839844-4.15625-2.34375-5.65625-1.5-1.503906-3.535156-2.34375-5.65625-2.34375zm-8 96h-24c-3.03125 0-5.800781 1.710938-7.15625 4.421875l-7.421875 14.835937-41.421875 41.429688-12.6875-12.6875 18.34375-18.34375c2.289062-2.289062 2.972656-5.730469 1.734375-8.71875s-4.15625-4.9375-7.390625-4.9375h-128v-16h80c4.417969 0 8-3.582031 8-8s-3.582031-8-8-8h-80v-16h80c4.417969 0 8-3.582031 8-8s-3.582031-8-8-8h-80v-16h80c4.417969 0 8-3.582031 8-8s-3.582031-8-8-8h-80v-16h157.578125l21.984375 14.65625c1.3125.875 2.859375 1.34375 4.4375 1.34375h24zm0 0"/>
                                    </svg>
                                </h1>
                            </div>
                            <h1 className="ml-2 text-xl font-extrabold">TestRelay</h1>
                        </div>


                        <div className="">
                            <ul className="flex flex-">
                                <li className="">
                                    <a className="flex items-center justify-center px-4 py-3 text-base font-medium rounded-md hover:text-gray-500"
                                       href="//app.testrelay.io/login">Login</a>
                                </li>
                                <li className="hidden sm:block px-2">
                                    <a className="shadow-md flex items-center justify-center px-4 py-2 text-base font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-700"
                                       href="//app.testrelay.io/register">Signup for free</a>
                                </li>
                            </ul>
                        </div>
                    </div>
                </div>

                <div className="mx-auto max-w-7xl sm:px-6 lg:px-8">
                    <div className="relative px-4 py-2 md:py-16 overflow-hidden sm:px-0">
                        <div>
                            <div className="flex flex-col items-center justify-center w-full h-full text-center ">
                                <h1 className="mt-4 text-4xl md:text-6xl text-gray-800 font-extrabold md:leading-none md:tracking-tight">
                                    Automate Your Take Home Assignments
                                </h1>


                                <h2 className="mt-5 text-lg text-gray-600 sm:mt-8 sm:max-w-2xl sm:mx-auto md:mt-8 md:text-xl lg:mx-0">
                                    Testrelay is an <span className="text-indigo-500">all-in-one</span> platform for
                                    managing your technical interviews through <span
                                    className="text-indigo-500">Github</span>.
                                </h2>
                                <div className="w-full mt-5 sm:mt-12 sm:flex items-center justify-center">
                                    <a className="block mx-auto sm:w-auto  shadow-md px-4 py-3 text-base font-medium text-white bg-indigo-600 rounded-md hover:bg-indigo-700"
                                       href="//app.testrelay.io/register">
                                        Get Started Now - It's free
                                    </a>
                                </div>
                                <div className=" w-full">
                                    <div className="max-w-3xl mx-auto my-10 shadow-md rounded-b-lg">
                                        <div
                                            className="w-full h-11 rounded-t-lg bg-gray-200 flex justify-start items-center space-x-1.5 px-3">
                                            <span className="w-3 h-3 rounded-full bg-red-400"/>
                                            <span className="w-3 h-3 rounded-full bg-yellow-400"/>
                                            <span className="w-3 h-3 rounded-full bg-green-400"/>
                                        </div>
                                        <div
                                            className="bg-gray-100 border-t-0 w-full  rounded-b-lg flex overflow-hidden"
                                            style={{height: "28rem"}}>
                                            <div class="pl-2 sm:pl-4 md:pl-14 text-left">
                                                <div className="relative wrap overflow-hidden h-full">
                                                    <div className="absolute border-gray-200 h-full border-2 ml-4"/>
                                                    <div className=" mt-4 flex items-center w-full">
                                                        <div
                                                            className="z-20 flex items-center order-1 bg-indigo-500 w-8 h-8 rounded-full">
                                                            <h1 className="mx-auto font-semibold text-lg text-white">
                                                                <svg xmlns="http://www.w3.org/2000/svg"
                                                                     className="h-5 w-5" fill="none" viewBox="0 0 24 24"
                                                                     stroke="currentColor">
                                                                    <path strokeLinecap="round" strokeLinejoin="round"
                                                                          strokeWidth={2}
                                                                          d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
                                                                </svg>
                                                            </h1>
                                                        </div>
                                                        <div
                                                            className="order-1 bg-white rounded-lg w-3/4 px-6 py-4 ml-4 sm:ml-8">
                                                            <h4 className="mb-1 sm:mb-3 text-gray-300 text-xs sm:text-sm">5
                                                                days ago</h4>
                                                            <p className="text-xs sm:text-sm leading-relaxed tracking-wide text-gray-700">TestRelay
                                                                sent <span
                                                                    className="text-indigo-500">jimmy@code.sh</span> the <span
                                                                    className="font-semibold">BE engineering test</span> with
                                                                a time limit of 3hrs. They have until the 6th of
                                                                September to respond.</p>
                                                        </div>
                                                    </div>
                                                    <div className=" mt-4 flex items-center w-full">
                                                        <div
                                                            className="z-20 flex items-center order-1 bg-indigo-500 w-8 h-8 rounded-full">
                                                            <h1 className="mx-auto font-semibold text-lg text-white">
                                                                <svg xmlns="http://www.w3.org/2000/svg"
                                                                     className="h-5 w-5" fill="none" viewBox="0 0 24 24"
                                                                     stroke="currentColor">
                                                                    <path strokeLinecap="round" strokeLinejoin="round"
                                                                          strokeWidth={2}
                                                                          d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                                                                </svg>
                                                            </h1>
                                                        </div>
                                                        <div
                                                            className="order-1 bg-white rounded-lg w-3/4 px-6 py-4 ml-4 sm:ml-8">
                                                            <h4 className="mb-1 sm:mb-3 text-gray-300 text-xs sm:text-sm">4
                                                                days ago</h4>
                                                            <p className="text-xs sm:text-sm leading-relaxed tracking-wide text-gray-700">
                                                                <span
                                                                    className="text-indigo-500">jimmy@code.sh</span> scheduled
                                                                a test for 4th September @ 4PM EST.</p>
                                                        </div>
                                                    </div>
                                                    <div className=" mt-4 flex items-center w-full">
                                                        <div
                                                            className="z-20 flex items-center order-1 bg-indigo-500 w-8 h-8 rounded-full">
                                                            <h1 className="mx-auto font-semibold text-lg text-white">
                                                                <svg xmlns="http://www.w3.org/2000/svg"
                                                                     className="h-5 w-5" fill="none" viewBox="0 0 24 24"
                                                                     stroke="currentColor">
                                                                    <path strokeLinecap="round" strokeLinejoin="round"
                                                                          strokeWidth={2}
                                                                          d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z"/>
                                                                </svg>
                                                            </h1>
                                                        </div>
                                                        <div
                                                            className="order-1 bg-white rounded-lg w-3/4 px-6 py-4 ml-4 sm:ml-8">
                                                            <h4 className="mb-1 sm:mb-3 text-gray-300 text-xs sm:text-sm">3
                                                                days ago @ 16:00</h4>
                                                            <p className="text-xs sm:text-sm leading-relaxed tracking-wide text-gray-700">TestRelay
                                                                has invited <span
                                                                    className="text-indigo-500">@jimmycode</span> to the
                                                                assignment repository <span
                                                                    className="text-indigo-500 underline">testrelay/jimmycode-xyx</span>.
                                                                Test is in progress.</p>
                                                        </div>
                                                    </div>
                                                    <div className=" mt-4 flex items-center w-full">
                                                        <div
                                                            className="z-20 flex items-center order-1 bg-green-500 w-8 h-8 rounded-full">
                                                            <h1 className="mx-auto font-semibold text-lg text-white">
                                                                <svg xmlns="http://www.w3.org/2000/svg"
                                                                     className="h-5 w-5" fill="none" viewBox="0 0 24 24"
                                                                     stroke="currentColor">
                                                                    <path strokeLinecap="round" strokeLinejoin="round"
                                                                          strokeWidth={2}
                                                                          d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                                                                </svg>
                                                            </h1>
                                                        </div>
                                                        <div
                                                            className="order-1 bg-white rounded-lg w-3/4 px-6 py-4 ml-4 sm:ml-8">
                                                            <h4 className="mb-1 sm:mb-3 text-gray-300 text-xs sm:text-sm">3
                                                                days ago @ 20:00</h4>
                                                            <p className="text-xs sm:text-sm leading-relaxed tracking-wide text-gray-700">
                                                                <span className="text-indigo-500">@jimmycode</span> has
                                                                submitted his assessment to the BE engineering test.
                                                                Code is ready for review.</p>
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
                            <div className="py-16 px-4 sm:px-8 text-center">
                                <svg className="h-20 w-20 mx-auto text-gray-800 group-hover:text-gray-200 mb-4"
                                     viewBox="0 0 256 250" version="1.1" preserveAspectRatio="xMidYMid">
                                    <g>
                                        <path fill="currentColor"
                                              d="M128.00106,0 C57.3172926,0 0,57.3066942 0,128.00106 C0,184.555281 36.6761997,232.535542 87.534937,249.460899 C93.9320223,250.645779 96.280588,246.684165 96.280588,243.303333 C96.280588,240.251045 96.1618878,230.167899 96.106777,219.472176 C60.4967585,227.215235 52.9826207,204.369712 52.9826207,204.369712 C47.1599584,189.574598 38.770408,185.640538 38.770408,185.640538 C27.1568785,177.696113 39.6458206,177.859325 39.6458206,177.859325 C52.4993419,178.762293 59.267365,191.04987 59.267365,191.04987 C70.6837675,210.618423 89.2115753,204.961093 96.5158685,201.690482 C97.6647155,193.417512 100.981959,187.77078 104.642583,184.574357 C76.211799,181.33766 46.324819,170.362144 46.324819,121.315702 C46.324819,107.340889 51.3250588,95.9223682 59.5132437,86.9583937 C58.1842268,83.7344152 53.8029229,70.715562 60.7532354,53.0843636 C60.7532354,53.0843636 71.5019501,49.6441813 95.9626412,66.2049595 C106.172967,63.368876 117.123047,61.9465949 128.00106,61.8978432 C138.879073,61.9465949 149.837632,63.368876 160.067033,66.2049595 C184.49805,49.6441813 195.231926,53.0843636 195.231926,53.0843636 C202.199197,70.715562 197.815773,83.7344152 196.486756,86.9583937 C204.694018,95.9223682 209.660343,107.340889 209.660343,121.315702 C209.660343,170.478725 179.716133,181.303747 151.213281,184.472614 C155.80443,188.444828 159.895342,196.234518 159.895342,208.176593 C159.895342,225.303317 159.746968,239.087361 159.746968,243.303333 C159.746968,246.709601 162.05102,250.70089 168.53925,249.443941 C219.370432,232.499507 256,184.536204 256,128.00106 C256,57.3066942 198.691187,0 128.00106,0 Z M47.9405593,182.340212 C47.6586465,182.976105 46.6581745,183.166873 45.7467277,182.730227 C44.8183235,182.312656 44.2968914,181.445722 44.5978808,180.80771 C44.8734344,180.152739 45.876026,179.97045 46.8023103,180.409216 C47.7328342,180.826786 48.2627451,181.702199 47.9405593,182.340212 Z M54.2367892,187.958254 C53.6263318,188.524199 52.4329723,188.261363 51.6232682,187.366874 C50.7860088,186.474504 50.6291553,185.281144 51.2480912,184.70672 C51.8776254,184.140775 53.0349512,184.405731 53.8743302,185.298101 C54.7115892,186.201069 54.8748019,187.38595 54.2367892,187.958254 Z M58.5562413,195.146347 C57.7719732,195.691096 56.4895886,195.180261 55.6968417,194.042013 C54.9125733,192.903764 54.9125733,191.538713 55.713799,190.991845 C56.5086651,190.444977 57.7719732,190.936735 58.5753181,192.066505 C59.3574669,193.22383 59.3574669,194.58888 58.5562413,195.146347 Z M65.8613592,203.471174 C65.1597571,204.244846 63.6654083,204.03712 62.5716717,202.981538 C61.4524999,201.94927 61.1409122,200.484596 61.8446341,199.710926 C62.5547146,198.935137 64.0575422,199.15346 65.1597571,200.200564 C66.2704506,201.230712 66.6095936,202.705984 65.8613592,203.471174 Z M75.3025151,206.281542 C74.9930474,207.284134 73.553809,207.739857 72.1039724,207.313809 C70.6562556,206.875043 69.7087748,205.700761 70.0012857,204.687571 C70.302275,203.678621 71.7478721,203.20382 73.2083069,203.659543 C74.6539041,204.09619 75.6035048,205.261994 75.3025151,206.281542 Z M86.046947,207.473627 C86.0829806,208.529209 84.8535871,209.404622 83.3316829,209.4237 C81.8013,209.457614 80.563428,208.603398 80.5464708,207.564772 C80.5464708,206.498591 81.7483088,205.631657 83.2786917,205.606221 C84.8005962,205.576546 86.046947,206.424403 86.046947,207.473627 Z M96.6021471,207.069023 C96.7844366,208.099171 95.7267341,209.156872 94.215428,209.438785 C92.7295577,209.710099 91.3539086,209.074206 91.1652603,208.052538 C90.9808515,206.996955 92.0576306,205.939253 93.5413813,205.66582 C95.054807,205.402984 96.4092596,206.021919 96.6021471,207.069023 Z"/>
                                    </g>
                                </svg>
                                <h2 className="font-bold text-4xl">TestRelay is completely <span
                                    className="text-indigo-600">open source</span></h2>
                                <p className="mt-2">Anyone can contribute to make TestRelay a better project.</p>
                                <div className="">
                                    <div
                                        className="mt-8 font-medium md:flex md:justify-center md:items-center space-y-3 md:space-y-0 md:space-x-5 md:flex-row">
                                        <div className="flex justify-center items-center">
                                            <div className="flex-shrink-0">
                                                <div
                                                    className="h-5 w-5 rounded-full bg-indigo-600 text-white flex items-center justify-center">
                                                    <svg xmlns="http://www.w3.org/2000/svg" className="h-3 w-3"
                                                         viewBox="0 0 20 20" fill="currentColor">
                                                        <path
                                                            d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z"/>
                                                    </svg>
                                                </div>
                                            </div>
                                            <p className="ml-2 leading-2 text-gray-600 text-left md:text-center">
                                                Star us on <a className="underline text-gray-500 hover:text-gray-600"
                                                              href="https://github.com/testrelay/testrelay">github</a>
                                            </p>
                                        </div>
                                        <div className="flex justify-center items-center">
                                            <div className="flex-shrink-0">
                                                <div
                                                    className="h-5 w-5 rounded-full bg-indigo-600 text-white flex items-center justify-center">
                                                    <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4"
                                                         fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                        <path strokeLinecap="round" strokeLinejoin="round"
                                                              strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
                                                    </svg>
                                                </div>
                                            </div>
                                            <p className="ml-2 leading-2 text-gray-600 text-left md:text-center">
                                                Add feature requests <a
                                                className="underline text-gray-500 hover:text-gray-600"
                                                href="https://github.com/testrelay/testrelay/issues">using issues</a>
                                            </p>
                                        </div>
                                        <div className="flex justify-center items-center text-left md:text-center">
                                            <div className="flex-shrink-0">
                                                <div
                                                    className="h-5 w-5 rounded-full bg-indigo-600 text-white flex items-center justify-center">
                                                    <svg role="img" viewBox="0 0 24 24" className="h-3 h-3"
                                                         xmlns="http://www.w3.org/2000/svg">
                                                        <path fill="currentColor"
                                                              d="M5.042 15.165a2.528 2.528 0 0 1-2.52 2.523A2.528 2.528 0 0 1 0 15.165a2.527 2.527 0 0 1 2.522-2.52h2.52v2.52zM6.313 15.165a2.527 2.527 0 0 1 2.521-2.52 2.527 2.527 0 0 1 2.521 2.52v6.313A2.528 2.528 0 0 1 8.834 24a2.528 2.528 0 0 1-2.521-2.522v-6.313zM8.834 5.042a2.528 2.528 0 0 1-2.521-2.52A2.528 2.528 0 0 1 8.834 0a2.528 2.528 0 0 1 2.521 2.522v2.52H8.834zM8.834 6.313a2.528 2.528 0 0 1 2.521 2.521 2.528 2.528 0 0 1-2.521 2.521H2.522A2.528 2.528 0 0 1 0 8.834a2.528 2.528 0 0 1 2.522-2.521h6.312zM18.956 8.834a2.528 2.528 0 0 1 2.522-2.521A2.528 2.528 0 0 1 24 8.834a2.528 2.528 0 0 1-2.522 2.521h-2.522V8.834zM17.688 8.834a2.528 2.528 0 0 1-2.523 2.521 2.527 2.527 0 0 1-2.52-2.521V2.522A2.527 2.527 0 0 1 15.165 0a2.528 2.528 0 0 1 2.523 2.522v6.312zM15.165 18.956a2.528 2.528 0 0 1 2.523 2.522A2.528 2.528 0 0 1 15.165 24a2.527 2.527 0 0 1-2.52-2.522v-2.522h2.52zM15.165 17.688a2.527 2.527 0 0 1-2.52-2.523 2.526 2.526 0 0 1 2.52-2.52h6.313A2.527 2.527 0 0 1 24 15.165a2.528 2.528 0 0 1-2.522 2.523h-6.313z"/>
                                                    </svg>
                                                </div>
                                            </div>
                                            <p className="ml-2 leading-2 text-gray-600 text-left md:text-center">
                                                Join the conversation <a
                                                className="underline text-gray-500 hover:text-gray-600"
                                                href="https://join.slack.com/t/newworkspace-up55834/shared_invite/zt-xtb6rzic-20b8K6yLT_trVgUEqnuYCQ">on
                                                slack</a>
                                            </p>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div
                        className="relative max-w-screen-xl p-4 px-4 py-16 sm:mt-10 mx-auto bg-white dark:bg-gray-800 sm:px-6 lg:px-8">
                        <div className="relative">
                            <div className="lg:grid lg:grid-flow-row-dense lg:grid-cols-2 lg:gap-8 lg:items-center">
                                <div className="ml-auto lg:ml-6 lg:col-start-2 lg:max-w-2xl text-center sm:text-left">
                                    <h4 className="mt-2 text-2xl font-extrabold leading-8 text-gray-900 dark:text-white sm:text-3xl sm:leading-9">
                                        Version Control Your Tests
                                    </h4>
                                    <p className="mt-4 text-lg leading-6 text-gray-600">
                                        Store your organisation's technical tests directly on github. Build your
                                        technical tests like you would applications.
                                    </p>

                                    <p className="mt-4 text-lg leading-6 text-gray-600">
                                        TestRelay automatically keeps up to date with the latest changes. Candidates
                                        always get the most recent version.
                                    </p>
                                </div>
                                <div className="bg-gray-50 relative px-4 py-8 sm:px-8 w-100 h-auto rounded">
                                    <div className="bg-white p-2 mb-8 overflow-hidden rounded shadow">
                                        <img src={Changelog}/>
                                    </div>

                                    <div
                                        className="bg-white w-2/3 sm:w-1/2 right-4 bottom-4 p-4 rounded shadow absolute">
                                        <label className="block text-gray-700 text-xs font-bold mb-2">Choose
                                            Candidate's Technical Assignment</label>
                                        <div className="w-full">
                                            <p className="mb-2 text-xs"> TestRelay will clone this repository to the
                                                candidate's technical test.</p>
                                            <div className="relative">
                                                <div className="bg-gray-100  rounded p-4 border"/>
                                                <div
                                                    className="grid text-xs divide-y shadow rounded absolute w-full bg-white top-2 left-2">
                                                    <div className="p-2 cursor-pointer">acme/be-eng-test</div>
                                                    <div
                                                        className="p-2 cursor-pointer text-white bg-indigo-500">acme/ds-interview
                                                        <img className="w-4 h-6 absolute z-10 right-6"
                                                             src={Cursor}/>
                                                    </div>
                                                    <div className="p-2 cursor-pointer">acme/fe-eng-test</div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div
                        className="relative max-w-screen-xl p-4 px-4 py-16 mx-auto bg-white dark:bg-gray-800 sm:px-6 lg:px-8">
                        <div className="relative">
                            <div className="lg:grid lg:grid-flow-row-dense lg:grid-cols-2 lg:gap-8 lg:items-center">
                                <div className="ml-auto lg:mr-6 lg:col-start-1 lg:max-w-2xl text-center sm:text-left">
                                    <h4 className="mt-2 text-2xl font-extrabold leading-8 text-gray-900 dark:text-white sm:text-3xl sm:leading-9">
                                        Automated Assignment Scheduling
                                    </h4>
                                    <p className="mt-4 text-lg leading-6 text-gray-600">
                                        Let TestRelay handle scheduling and running a candidate's technical test.
                                        Candidates receive a dedicated portal to schedule tests when it suits them.
                                    </p>

                                    <p className="mt-4 text-lg leading-6 text-gray-600">
                                        TestRelay automatically keeps on top of candidates technical assignments,
                                        managing access for candidates when needed. Never worry about sending technical
                                        tests on your days off.
                                    </p>
                                </div>
                                <div className="relative mt-10 space-y-6 relative-20 lg:mt-0">
                                    <div className="bg-gray-50 relative px-4 py-8 sm:px-8 w-100 h-auto rounded">
                                        <div>
                                            <div className="bg-white mb-2 p-2 overflow-hidden rounded shadow">
                                                <div
                                                    className="relative bg-white p-2 rounded text-center md:text-left">
                                                    <div className="grid grid-cols-3 gap-4">
                                                        <div>
                                                            <div
                                                                className="text-xs font-medium text-indigo-700 mb-1">
                                                                Acme Org Test
                                                            </div>
                                                            <div
                                                                className="text-xs text-gray-500 flex items-center justify-center md:justify-start">
                                                                <svg xmlns="http://www.w3.org/2000/svg"
                                                                     className="h-4 w-4 flex-shrink-0" fill="none"
                                                                     viewBox="0 0 24 24" stroke="currentColor">
                                                                    <path strokeLinecap="round"
                                                                          strokeLinejoin="round" strokeWidth={2}
                                                                          d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                                                                </svg>
                                                                <span className="ml-1">3 hours</span>
                                                            </div>
                                                        </div>
                                                        <div>
                                                            <div>
                                                                <div
                                                                    className="text-xs mb-1 font-medium text-gray-900">
                                                                    8th March 2021
                                                                </div>
                                                                <div className="text-xs text-gray-500">
                                                                    at 15:00 GMT
                                                                </div>
                                                            </div>
                                                        </div>
                                                        <div className="flex justify-end items-center">
                                                            <svg xmlns="http://www.w3.org/2000/svg"
                                                                 className="h-4 w-4" fill="none" viewBox="0 0 24 24"
                                                                 stroke="currentColor">
                                                                <path strokeLinecap="round" strokeLinejoin="round"
                                                                      strokeWidth={2} d="M9 5l7 7-7 7"/>
                                                            </svg>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                            <div className="bg-white mb-2 p-2 overflow-hidden rounded shadow">
                                                <div
                                                    className="relative bg-white p-2 rounded text-center md:text-left">
                                                    <div className="grid grid-cols-3 gap-4">
                                                        <div>
                                                            <div
                                                                className="text-xs font-medium text-indigo-700 mb-1">
                                                                Beta Labs Test
                                                            </div>
                                                            <div
                                                                className="text-xs text-gray-500 flex items-center justify-center md:justify-start">
                                                                <svg xmlns="http://www.w3.org/2000/svg"
                                                                     className="h-4 w-4 flex-shrink-0" fill="none"
                                                                     viewBox="0 0 24 24" stroke="currentColor">
                                                                    <path strokeLinecap="round"
                                                                          strokeLinejoin="round" strokeWidth={2}
                                                                          d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                                                                </svg>
                                                                <span className="ml-1">45 minutes</span>
                                                            </div>
                                                        </div>
                                                        <div>
                                                            <div>
                                                                <div
                                                                    className="text-xs mb-1 font-medium text-gray-900">
                                                                    15th March 2021
                                                                </div>
                                                                <div className="text-xs text-gray-500">
                                                                    at 09:00 GMT
                                                                </div>
                                                            </div>
                                                        </div>
                                                        <div className="flex justify-end items-center">
                                                            <svg xmlns="http://www.w3.org/2000/svg"
                                                                 className="h-4 w-4" fill="none" viewBox="0 0 24 24"
                                                                 stroke="currentColor">
                                                                <path strokeLinecap="round" strokeLinejoin="round"
                                                                      strokeWidth={2} d="M9 5l7 7-7 7"/>
                                                            </svg>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                            <div className="bg-white mb-8 p-2 overflow-hidden rounded shadow">
                                                <div
                                                    className="relative bg-white p-2 rounded text-center md:text-left">
                                                    <div className="grid grid-cols-2 gap-4">
                                                        <div>
                                                            <div
                                                                className="text-xs font-medium text-indigo-700 mb-1">
                                                                Theta AI Test
                                                            </div>
                                                            <div
                                                                className="text-xs text-gray-500 flex items-center justify-center md:justify-start">
                                                                <svg xmlns="http://www.w3.org/2000/svg"
                                                                     className="h-4 w-4 flex-shrink-0" fill="none"
                                                                     viewBox="0 0 24 24" stroke="currentColor">
                                                                    <path strokeLinecap="round"
                                                                          strokeLinejoin="round" strokeWidth={2}
                                                                          d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                                                                </svg>
                                                                <span className="ml-1">1 hour</span>
                                                            </div>
                                                        </div>
                                                        <div className="text-right">
                                                            <div>
                                                                <div
                                                                    className="text-xs mb-1 font-medium text-green-600">
                                                                    Completed
                                                                </div>
                                                                <div className="text-xs text-gray-500">
                                                                    8th March 2021 @ 15:00 GMT
                                                                </div>
                                                            </div>
                                                        </div>
                                                    </div>
                                                </div>
                                            </div>
                                            <div
                                                className="absolute bg-gray-50  opacity-30 w-full h-full right-0 top-0 bottom-0"/>
                                            <div
                                                className="bg-white w-3/4 sm:w-1/2 left-0 right-0 top-20 mx-auto p-4 pb-10 rounded shadow-xl border  absolute">
                                                <label className="block text-gray-700 text-xs font-bold mb-2">Nipo
                                                    Inc has invited you to a technical test</label>
                                                <div className="w-full h-auto relative">
                                                    <p className="text-xs mb-2">Schedule a date & time for your
                                                        assignment.</p>
                                                    <div
                                                        className="px-2 py-2 text-xs border bg-gray-50 text-gray-800 flex rounded items-center mb-2">
                                                        <span className="flex-grow">Tuesday 12th, Oct 2021</span>
                                                        <div className="text-right">
                                                            <svg xmlns="http://www.w3.org/2000/svg"
                                                                 className="h-4 w-4" fill="none" viewBox="0 0 24 24"
                                                                 stroke="currentColor">
                                                                <path strokeLinecap="round" strokeLinejoin="round"
                                                                      strokeWidth={2}
                                                                      d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                                                            </svg>
                                                        </div>
                                                    </div>
                                                    <div
                                                        className="shadow-xl border p-2 bg-white absolute z-10 rounded w-3/4 sm:w-2/3">
                                                        <img className="w-full" src={DatePicker}/>
                                                    </div>
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div
                        className="relative max-w-screen-xl p-4 px-4 pb-10 pt-24 sm:pb-16 sm:pt-16 mx-auto bg-white dark:bg-gray-800 sm:px-6 lg:px-8">
                        <div className="relative">
                            <div
                                className="lg:grid lg:grid-flow-row-dense lg:grid-cols-2 lg:gap-8 lg:items-center text-center sm:text-left">
                                <div className="ml-auto lg:ml-6 lg:col-start-2 lg:max-w-2xl">
                                    <h4 className="mt-2 text-2xl font-extrabold leading-8 text-gray-900 dark:text-white sm:text-3xl sm:leading-9">
                                        Secure Test Environments
                                    </h4>
                                    <p className="mt-4 text-lg leading-6 text-gray-600">
                                        TestRelay auto generates private repositories that candidates will use to sit
                                        their technical test.
                                    </p>
                                    <p className="mt-4 text-lg leading-6 text-gray-600">
                                        Once a candidate finishes their assignment, they're removed from the
                                        repository.
                                    </p>
                                </div>
                                <div className="bg-gray-50 px-4 py-8 sm:px-8 rounded">
                                    <img className="mx-auto py-4 w-20"
                                         src="https://ci3.googleusercontent.com/proxy/UFMFbT0UrmdnumWhFvGkEviupuQf7QIzl3VYfzne8CB_LgHu7d0dAbUn5An4x3f8ySn_6KomW968hn2Z6wR3gk3e3r0L-r8TYF530A4AASxTI0c=s0-d-e1-ft#https://github.githubassets.com/images/email/global/wordmark.png"/>
                                    <div className="bg-white border rounded p-4 w-full sm:w-5/6 mx-auto">
                                        <div className="py-2 grid grid-cols-3 mx-auto w-1/2 sm:w-1/3">
                                            <div className="flex items-center text-center">
                                                <img className="mx-auto w-10 h-10"
                                                     src="https://avatars.githubusercontent.com/u/90466910?s=120"
                                                     alt="@testrelay"/>
                                            </div>
                                            <div className="flex items-center text-center">
                                                <img className="mx-auto w-10 h-10" alt="plus"
                                                     src="https://github.githubassets.com/images/email/organization/octicon-plus.png"/>
                                            </div>
                                            <div className="flex items-center text-center">
                                                <img className="mx-auto w-10 h-10"
                                                     src="https://avatars.githubusercontent.com/u/6455139?s=120"
                                                     alt="@hugorut"/>
                                            </div>
                                        </div>
                                        <div>
                                            <p className="text-center text-sm mt-2 mb-4">@testrelay has invited you
                                                to collaborate on the<br/><span
                                                    className="font-bold">testrelay/hugorut-nipo-inc-test-hyly</span> repository
                                            </p>
                                        </div>
                                        <div className="border-b mb-2"/>
                                        <div>
                                            <p className="text-xs mb-4">You can <span className="text-blue-500">accept or decline</span> this
                                                invitation. You can also visit <span
                                                    className="text-blue-500">@testrelay</span> to learn a bit more
                                                about them.</p>
                                            <p className="text-xs">This invitation will expire in 7 days.</p>
                                            <div
                                                className="cursor-pointer bg-gradient-to-b  font-bold from-green-400 to-green-600 text-center w-32  mx-auto text-xs rounded px-4 py-2 mt-6 text-white">View
                                                Invitation
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>

                    <div
                        className="relative max-w-screen-xl p-4 px-4 pb-16 sm:pt-16 mx-auto bg-white dark:bg-gray-800 sm:px-6 lg:px-8">
                        <div className="relative">
                            <div
                                className="lg:grid lg:grid-flow-row-dense lg:grid-cols-2 lg:gap-8 lg:items-center text-center sm:text-left">
                                <div className="ml-auto lg:mr-6 lg:max-w-2xl">
                                    <h4 className="mt-2 text-2xl font-extrabold leading-8 text-gray-900 dark:text-white sm:text-3xl sm:leading-9">
                                        Evaluate Candidates Collaboratively
                                    </h4>
                                    <p className="mt-4 text-lg leading-6 text-gray-600">
                                        Import colleagues into TestRelay and seamlessly manage which assignments they
                                        need to review.
                                    </p>
                                    <p className="mt-4 text-lg leading-6 text-gray-600">
                                        Completed technical assignments are opened as PRs on Github. TestRelay auto
                                        requests a code review from your nominated interviewers.
                                    </p>

                                </div>
                                <div className="bg-gray-50 relative px-4 py-8 sm:px-8  rounded overflow-hidden">
                                    <div className="mx-auto sm:hidden h-96 w-auto bg-left bg-no-repeat" style={{
                                        "backgroundImage": "url(" + Review + ")",
                                        "backgroundSize": "220%",
                                        "backgroundPositionX": "-86px"
                                    }}/>
                                    <div className="hidden sm:block bg-white p-2 overflow-hidden rounded shadow">
                                        <img src={Review}/>
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
                                    <p className="lg:w-2/3 mx-auto leading-relaxed text-base text-gray-600">We're
                                        committed to having a free plan, forever. Only pay for features when you can
                                        afford to do so.</p>
                                </div>
                                <div className="flex flex-wrap justify-center -m-4">
                                    <div className="p-4 xl:w-1/4 md:w-1/2 w-full">
                                        <div
                                            className="h-full p-6 rounded-lg bg-white  border-2 border-indigo-500 shadow-md flex flex-col relative overflow-hidden">
                                            <h2 className="text-sm tracking-widest title-font mb-1 font-medium">START</h2>
                                            <h1 className="text-5xl text-gray-900 pb-4 mb-4 border-b border-gray-200 leading-none">Free</h1>
                                            <p className="flex items-center text-gray-600 mb-2">
                        <span
                            className="w-4 h-4 mr-2 inline-flex items-center justify-center bg-indigo-600 text-white rounded-full flex-shrink-0">
                          <svg fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                               stroke-width="2.5" className="w-3 h-3" viewBox="0 0 24 24">
                            <path d="M20 6L9 17l-5-5"></path>
                          </svg>
                        </span>Unlimited assignments
                                            </p>
                                            <p className="flex items-center text-gray-600 mb-2">
                        <span
                            className="w-4 h-4 mr-2 inline-flex items-center justify-center bg-indigo-600 text-white rounded-full flex-shrink-0">
                          <svg fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                               stroke-width="2.5" className="w-3 h-3" viewBox="0 0 24 24">
                            <path d="M20 6L9 17l-5-5"></path>
                          </svg>
                        </span>Unlimited users
                                            </p>
                                            <p className="flex items-center text-gray-600 mb-6">
                        <span
                            className="w-4 h-4 mr-2 inline-flex items-center justify-center bg-indigo-600 text-white rounded-full flex-shrink-0">
                          <svg fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round"
                               stroke-width="2.5" className="w-3 h-3" viewBox="0 0 24 24">
                            <path d="M20 6L9 17l-5-5"></path>
                          </svg>
                        </span>Everything for free!
                                            </p>
                                            <a href="//app.testrelay.com/register"
                                               className="flex items-center mt-auto text-white bg-indigo-500 bborder-0 py-2 px-4 w-full focus:outline-none hover:bg-gray-500 rounded">Get
                                                Started
                                                <svg fill="none" stroke="currentColor" stroke-linecap="round"
                                                     stroke-linejoin="round" stroke-width="2"
                                                     className="w-4 h-4 ml-auto" viewBox="0 0 24 24">
                                                    <path d="M5 12h14M12 5l7 7-7 7"></path>
                                                </svg>
                                            </a>
                                            <p className="text-xs text-gray-500 mt-3">What are you waiting for? It's
                                                free!</p>
                                        </div>
                                    </div>
                                    <div className="p-4 xl:w-1/4 md:w-1/2 w-full">
                                        <div
                                            className="h-full p-6 rounded-lg  bg-white shadow-md flex flex-col relative overflow-hidden">
                                            <span
                                                className="bg-gray-300 text-white px-3 py-1 tracking-widest text-xs absolute right-0 top-0 rounded-bl">coming when it's ready</span>
                                            <h2 className="text-sm tracking-widest title-font mb-1 font-medium">PRO</h2>
                                            <h1 className="text-5xl text-gray-900 leading-none flex items-center pb-4 mb-4 border-b border-gray-200">
                                                <span>$?</span>
                                                <span className="text-lg ml-1 font-normal text-gray-500">/mo</span>
                                            </h1>
                                            <p className="flex items-center text-gray-600 mb-2">
                                                We're focused on building a product that people will love. All plans are
                                                free for now and we will always have a free tier.
                                            </p>
                                            <p className="text-xs text-gray-500 mt-3">Want a say in what's in the pro
                                                plan, email <span className="text-indigo-500">hello@testrelay.io.</span>
                                            </p>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </section>
                    </div>
                </div>
            </main>
            <footer className="text-gray-600 body-font">
                <div className="container px-5 py-8 mx-auto flex items-center sm:flex-row flex-col">
                    <a className="flex title-font font-medium items-center md:justify-start justify-center text-gray-900">
                        <div className="flex items-center bg-gray-800 shadow w-10 h-10 rounded-full">
                            <h1 className="mx-auto font-semibold text-lg text-white">
                                <svg height="1.8rem" viewBox="0 -48 480 480" width="1.8rem" fill="#fff"
                                     xmlns="http://www.w3.org/2000/svg">
                                    <path
                                        d="m232 0h-64c-3.617188.00390625-6.785156 2.429688-7.726562 5.921875l-45.789063 170.078125h-50.484375c-3.441406 0-6.5 2.203125-7.585938 5.46875l-14.179687 42.53125h-34.234375c-2.121094 0-4.15625.839844-5.65625 2.34375-1.503906 1.5-2.34375 3.535156-2.34375 5.65625v112c0 2.121094.839844 4.15625 2.34375 5.65625 1.5 1.503906 3.535156 2.34375 5.65625 2.34375h69.757812l-5.515624 22.0625c-.597657 2.390625-.0625 4.921875 1.453124 6.859375 1.515626 1.941406 3.84375 3.078125 6.304688 3.078125h64c3.671875 0 6.871094-2.5 7.757812-6.0625l6.484376-25.9375h9.757812c3.8125-.003906 7.09375-2.691406 7.84375-6.429688l14.710938-73.570312h9.445312c2.121094 0 4.15625-.839844 5.65625-2.34375 1.503906-1.5 2.34375-3.535156 2.34375-5.65625v-48c0-2.121094-.839844-4.15625-2.34375-5.65625-1.5-1.503906-3.535156-2.34375-5.65625-2.34375h-13.558594l53.285156-197.921875c.648438-2.402344.140626-4.972656-1.375-6.945313-1.515624-1.976562-3.863281-3.132812-6.351562-3.132812zm-94.25 368h-47.5l4-16h47.5zm54.25-112h-8c-3.8125.003906-7.09375 2.691406-7.84375 6.429688l-14.710938 73.570312h-145.445312v-96h32c3.441406 0 6.5-2.203125 7.585938-5.46875l14.179687-42.53125h40.410156l-5.902343 21.921875c-.648438 2.402344-.140626 4.972656 1.375 6.945313 1.515624 1.976562 3.863281 3.132812 6.351562 3.132812h80zm-22.132812-48h-47.429688l51.695312-192h47.429688zm0 0"/>
                                    <path
                                        d="m472 208h-29.578125l-21.984375-14.65625c-1.3125-.875-2.859375-1.34375-4.4375-1.34375h-168c-2.121094 0-4.15625.839844-5.65625 2.34375-1.503906 1.5-2.34375 3.535156-2.34375 5.65625v128c0 2.121094.839844 4.15625 2.34375 5.65625 1.5 1.503906 3.535156 2.34375 5.65625 2.34375h116.6875l-10.34375 10.34375c-3.125 3.125-3.125 8.1875 0 11.3125l24 24c3.125 3.125 8.1875 3.125 11.3125 0l48-48c.609375-.609375 1.113281-1.308594 1.5-2.078125l5.789062-11.578125h27.054688c2.121094 0 4.15625-.839844 5.65625-2.34375 1.503906-1.5 2.34375-3.535156 2.34375-5.65625v-96c0-2.121094-.839844-4.15625-2.34375-5.65625-1.5-1.503906-3.535156-2.34375-5.65625-2.34375zm-8 96h-24c-3.03125 0-5.800781 1.710938-7.15625 4.421875l-7.421875 14.835937-41.421875 41.429688-12.6875-12.6875 18.34375-18.34375c2.289062-2.289062 2.972656-5.730469 1.734375-8.71875s-4.15625-4.9375-7.390625-4.9375h-128v-16h80c4.417969 0 8-3.582031 8-8s-3.582031-8-8-8h-80v-16h80c4.417969 0 8-3.582031 8-8s-3.582031-8-8-8h-80v-16h80c4.417969 0 8-3.582031 8-8s-3.582031-8-8-8h-80v-16h157.578125l21.984375 14.65625c1.3125.875 2.859375 1.34375 4.4375 1.34375h24zm0 0"/>
                                </svg>
                            </h1>
                        </div>
                        <span className="ml-3 text-xl">TestRelay</span>
                    </a>
                    <p className="text-sm text-gray-500 sm:ml-4 sm:pl-4 sm:border-l-2 sm:border-gray-200 sm:py-2 sm:mt-0 mt-4">©
                        2021 TestRelay</p>
                </div>
            </footer>
        </div>
    )
}

export default IndexPage
