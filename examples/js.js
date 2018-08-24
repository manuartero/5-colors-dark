import test from 'ava';
import {withBrowser} from '../test-utils/e2e/browser';

test(
    'topup invoice flow',
    withBrowser(async (t, {openPage, select, server}) => {
        server.stubApi({
            call: ['Auth', 'getSessionInfo'],
            returns: server.sessionInfo,
        });

        server.stubApi({
            call: ['AccountDashboard', 'getAccountBalance'],
            returns: {
                isAvailable: true,
                moneyAmount: {
                    amount: 1200,
                    precision: 2,
                    currency: 'ARS',
                },
            },
        });

        server.stubApi({
            call: ['AccountTopups', 'getTopupInvoiceData'],
            returns: {
                success: true,
                response: {
                    minAmount: {
                        amount: 100,
                        currency: 'ARS',
                        precision: 2,
                    },
                    maxAmount: {
                        amount: 15000,
                        currency: 'ARS',
                        precision: 2,
                    },
                    termsUrl: 'http://whatever.com',
                },
            },
        });

        server.stubApi({
            call: ['AccountSubscriptions', 'getSubscriptions'],
            returns: {
                success: true,
                response: [
                    {id: '12345679', msisdn: '34677777777', lifecycleStatus: 'active', isSelected: true},
                ],
            },
        });

        /* TODO MOVAR-1301: restore the real test */
        // const requestTopupMock = server.mockApi({
        //     call: ['AccountFlow', 'execute'],
        //     returns: {success: true},
        // });
        /* --- */

        const page = await openPage({path: '/pages/topup-invoice', withSession: server.sessionInfo});

        t.true(await page.hasElement(select('disabled-mode')), 'shows disabled mode'); // MOVAR-1301

        /* TODO MOVAR-1301: restore the real test */
        // await page.type(select('amount-input'), '20');
        // await page.click(select('confirm-button'));
        // await page.click(select('topup-button'));
        // await waitForCondition(requestTopupMock.calledOnce, 'Flow api is called');
        // t.pass('Topup done');
        /* --- */
    })
);
