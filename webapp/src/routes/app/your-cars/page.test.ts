import { expect, test, describe, beforeAll, afterAll, afterEach, vi } from 'vitest';
import { http, HttpResponse } from 'msw';
import { load } from './+page';
import { render, screen } from '@testing-library/svelte';
import { setupServer } from 'msw/node';
import { BACKEND_SERVER } from '$lib/constants';
import CarPage from './+page.svelte';
import { test_data } from './test_data';
import userEvent from '@testing-library/user-event';
import type { PageData, PageLoadEvent } from './$types';
import { mock } from 'vitest-mock-extended';

const server = setupServer();
const user = userEvent.setup();

beforeAll(() => {
    // NOTE: server.listen must be called before `createClient` is used to ensure
    // the msw can inject its version of `fetch` to intercept the requests.
    server.listen({
        onUnhandledRequest: (request) => {
            throw new Error(`No request handler found for ${request.method} ${request.url}`);
        }
    });
});


afterEach(() => server.resetHandlers());

afterAll(() => server.close());

describe('fetch cars information test', () => {
    test('test if cars is loaded correctly', async () => {
        //Data is loaded correctly
        server.use(
            http.get(`${BACKEND_SERVER}/cars`, () => HttpResponse.json(test_data, { status: 200 }))
        );

        const loadEvent = mock<PageLoadEvent>({ fetch: global.fetch });
        const data = (await load(loadEvent)) as PageData;
        expect(data.cars).toStrictEqual(test_data);

        render(CarPage, { data: data });
        test_data.forEach((car) => {
            screen.getByText(car.details.license_plate);
            screen.getByText('Color: ' + car.details.color);
            screen.getByText('Make: ' + car.details.make);
            screen.getByText('Model: ' + car.details.model);
        });
    });

    test('test if cars create work correctly', async () => {
        const data = mock<PageData>({ cars: test_data });
        const newCarDetail = {
            license_plate: 'lic-new',
            make: 'color-new',
            model: 'model-new',
            color: 'make-new'
        };
        render(CarPage, { data: data });

        server.use(
            http.post(`${BACKEND_SERVER}/cars`, () =>
                HttpResponse.json({ details: newCarDetail, id: 'random-id' }, { status: 201 })
            )
        );

        const createButton = screen.getByRole('button', { name: 'new-car-button' });
        await user.click(createButton);
        let licensePlateField = screen.getByRole('textbox', {
            name: 'License plate'
        });
        let colorField = screen.getByRole('textbox', { name: 'Color' });
        let modelField = screen.getByRole('textbox', { name: 'Model' });
        let makeField = screen.getByRole('textbox', { name: 'Make' });

        console.log('Test creating a new car with correct input');
        await user.click(licensePlateField);
        await user.keyboard(newCarDetail.license_plate);

        await user.click(colorField);
        await user.keyboard(newCarDetail.color);

        await user.click(modelField);
        await user.keyboard(newCarDetail.model);

        await user.click(makeField);
        await user.keyboard(newCarDetail.make);

        let confirmButton = screen.getByRole('button', { name: 'Confirm' });
        await user.click(confirmButton);
        screen.getByText(newCarDetail.license_plate);
        screen.getByText('Color: ' + newCarDetail.color);
        screen.getByText('Make: ' + newCarDetail.make);
        screen.getByText('Model: ' + newCarDetail.model);

        server.use(
            http.post(`${BACKEND_SERVER}/cars`, () =>
                HttpResponse.json(
                    {
                        details: 'the specified license plate is invalid',
                        errors: [
                            {
                                location: 'body.license_plate',
                                value: 'very wrong license'
                            }
                        ]
                    },
                    { status: 422 }
                )
            )
        );
        console.log('Test creating a new car with incorrect input');
        await user.click(createButton);

        licensePlateField = screen.getByRole('textbox', {
            name: 'License plate'
        });
        colorField = screen.getByRole('textbox', { name: 'Color' });
        modelField = screen.getByRole('textbox', { name: 'Model' });
        makeField = screen.getByRole('textbox', { name: 'Make' });
        confirmButton = screen.getByRole('button', { name: 'Confirm' });
        await user.click(licensePlateField);
        await user.keyboard('Very wrong license this is certainly invalid');

        await user.click(colorField);
        await user.keyboard('wrong');

        await user.click(modelField);
        await user.keyboard('wrong');

        await user.click(makeField);
        await user.keyboard('wrong');

        await user.click(confirmButton);
        screen.getByText('Error:');
        const closeModalButton = screen.getByRole('button', { name: 'Close the modal' });
        await user.click(closeModalButton);
    });

    test('Test delete functionality', async () => {
        const data = mock<PageData>({ cars: test_data });
        window.confirm = vi.fn(() => {
            console.log('confirm');
            return true;
        });

        render(CarPage, { data: data });

        server.use(
            http.delete(`${BACKEND_SERVER}/cars/:id`, () =>
                HttpResponse.json({ data: 'random' }, { status: 204 })
            )
        );
        const carToDelete = screen.getByText(test_data[0].details.license_plate);
        await user.click(carToDelete);
        const deleteButton = screen.getByRole('button', { name: 'Delete' });
        await user.click(deleteButton);
        const deletedCar = screen.queryByText(test_data[0].details.license_plate);
        expect(deletedCar).toBe(null);
    });

    test('Test edit functionality', async () => {
        const data = mock<PageData>({ cars: test_data });
        render(CarPage, { data: data });
        const car_to_edit = screen.getByText(test_data[0].details.license_plate);
        await user.click(car_to_edit);
        let editButton = screen.getByRole('button', { name: 'Edit' });
        await user.click(editButton);

        const editCarDetail = {
            license_plate: 'lic-edit',
            make: 'color-edit',
            model: 'model-edit',
            color: 'make-edit'
        };

        server.use(
            http.put(`${BACKEND_SERVER}/cars/:id`, () =>
                HttpResponse.json({ details: editCarDetail, id: test_data[0].id }, { status: 200 })
            )
        );
        console.log('Edit with correct input');
        let licensePlateField = screen.getByRole('textbox', { name: 'License plate' });
        let colorField = screen.getByRole('textbox', { name: 'Color' });
        let modelField = screen.getByRole('textbox', { name: 'Model' });
        let makeField = screen.getByRole('textbox', { name: 'Make' });
        await user.click(licensePlateField);
        await user.clear(licensePlateField);
        await user.keyboard(editCarDetail.license_plate);

        await user.click(colorField);
        await user.clear(colorField);
        await user.keyboard(editCarDetail.color);

        await user.click(modelField);
        await user.clear(modelField);
        await user.keyboard(editCarDetail.model);

        await user.click(makeField);
        await user.clear(makeField);
        await user.keyboard(editCarDetail.make);

        let confirmButton = screen.getByRole('button', { name: 'Confirm' });
        await user.click(confirmButton);
        screen.getAllByText(editCarDetail.license_plate);
        screen.getAllByText('Color: ' + editCarDetail.color);
        screen.getAllByText('Make: ' + editCarDetail.make);
        screen.getAllByText('Model: ' + editCarDetail.model);

        console.log('Edit with wrong input');
        server.use(
            http.put(`${BACKEND_SERVER}/cars/:id`, () =>
                HttpResponse.json(
                    {
                        details: 'the specified license plate is invalid',
                        errors: [
                            {
                                location: 'body.license_plate',
                                value: 'very wrong license'
                            }
                        ]
                    },
                    { status: 422 }
                )
            )
        );

        editButton = screen.getByRole('button', { name: 'Edit' });
        await user.click(editButton);
        licensePlateField = screen.getByRole('textbox', { name: 'License plate' });
        colorField = screen.getByRole('textbox', { name: 'Color' });
        modelField = screen.getByRole('textbox', { name: 'Model' });
        makeField = screen.getByRole('textbox', { name: 'Make' });
        await user.click(licensePlateField);
        await user.clear(licensePlateField);
        await user.keyboard('1');

        await user.click(colorField);
        await user.clear(colorField);
        await user.keyboard('1');

        await user.click(modelField);
        await user.clear(modelField);
        await user.keyboard('1');

        await user.click(makeField);
        await user.clear(makeField);
        await user.keyboard('1');
        confirmButton = screen.getByRole('button', { name: 'Confirm' });

        await user.click(confirmButton);

        screen.getByText('Error:');
        const closeModalButton = screen.getByRole('button', { name: 'Close' });
        await user.click(closeModalButton);
    });
});
